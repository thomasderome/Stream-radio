package services

import (
	"context"
	"io"
	"log"
	"net/http"
	"os/exec"
	"radio_stream/model"
	"sync"

	"github.com/ebitengine/oto/v3"
)

type PlayersServiceStruct struct {
	isPlaying bool
	volume    float64

	stationId   string
	stationUrl  string
	stationName string
	stationImg  string

	contextOto *oto.Context
	player     *oto.Player
	stopCtx    context.CancelFunc
	mut        sync.RWMutex
}

var PlayersService *PlayersServiceStruct

func InitPlayersService() {
	options := &oto.NewContextOptions{
		SampleRate:   44100,
		ChannelCount: 1,
		Format:       oto.FormatSignedInt16LE,
	}

	contextOto, readyChan, err := oto.NewContext(options)
	if err != nil {
		log.Fatal("Impossible to create context oto: ", err)
	}
	<-readyChan

	dataPlayer, err := GetStationData()

	PlayersService = new(PlayersServiceStruct{
		isPlaying:   dataPlayer.IsPlaying,
		volume:      dataPlayer.Volume,
		stationName: dataPlayer.StationName,
		stationUrl:  dataPlayer.StationUrl,
		stationImg:  dataPlayer.StationImg,
		contextOto:  contextOto,
	})

	if PlayersService.isPlaying && PlayersService.stationUrl != "" {
		err = PlayersService.playStream()
		if err != nil {
			log.Fatal("Impossible start service Player: ", err)
		}

		PlayersService.player.SetVolume(PlayersService.volume / 100)
		PlayersService.player.Play()
		PlayersService.updateIsPlaying()
	}
}

func (PlayersService *PlayersServiceStruct) updateDBAndPlayerService(stationId string) (err error) {
	_, err = DB.Exec("UPDATE state SET station_id = ? WHERE id=1", stationId)
	if err != nil {
		log.Printf("Impossible to update table state with link play: %s", err)
		return err
	}

	PlayersService.stationId = stationId

	stationdata, err := GetStationData()
	if err != nil {
		return err
	}

	PlayersService.stationUrl = stationdata.StationUrl
	PlayersService.stationName = stationdata.StationName
	PlayersService.stationImg = stationdata.StationImg

	return nil
}

func convertToPCM(input io.Reader) (io.ReadCloser, *exec.Cmd, error) {
	cmd := exec.Command("ffmpeg",
		"-i", "pipe:0",
		"-f", "s16le",
		"-ar", "44100",
		"-ac", "1",
		"pipe:1",
	)

	cmd.Stdin = input
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, nil, err
	}

	if err = cmd.Start(); err != nil {
		return nil, nil, err
	}

	return stdout, cmd, nil
}

func (PlayersService *PlayersServiceStruct) createPlayer(audio io.ReadCloser) {
	player := PlayersService.contextOto.NewPlayer(audio)

	PlayersService.mut.Lock()
	defer PlayersService.mut.Unlock()

	PlayersService.player = player
}

func (PlayersService *PlayersServiceStruct) playStream() (err error) {
	streamUrl := PlayersService.stationUrl
	resp, err := http.Get(streamUrl)
	if err != nil {
		log.Printf("Impossible to get url steam: %s", streamUrl)
		return err
	}

	audio, cmd, err := convertToPCM(resp.Body)
	if err != nil {
		resp.Body.Close()
		log.Printf("Impossible to extract audio in boddy: %s", err)
		return err
	}

	PlayersService.createPlayer(audio)
	PlayersService.player.Play()
	PlayersService.updateIsPlaying()

	ctx, cancel := context.WithCancel(context.Background())

	PlayersService.mut.Lock()
	PlayersService.stopCtx = cancel
	PlayersService.mut.Unlock()

	go keepAlive(ctx, resp, PlayersService.player, audio, cmd)

	return nil
}

func keepAlive(ctx context.Context, resp *http.Response, player *oto.Player, audio io.ReadCloser, cmd *exec.Cmd) {
	<-ctx.Done()

	_ = player.Close()
	_ = audio.Close()
	_ = resp.Body.Close()

	if cmd != nil {
		_ = cmd.Process.Kill()
		_ = cmd.Wait()
	}

	PlayersService.mut.Lock()
	defer PlayersService.mut.Unlock()

	if PlayersService.player != player {
		return
	}

	PlayersService.isPlaying = false
	PlayersService.player = nil
	PlayersService.stopCtx = nil
}

func (PlayersService *PlayersServiceStruct) updateIsPlaying() {
	PlayersService.mut.Lock()
	defer PlayersService.mut.Unlock()
	PlayersService.isPlaying = PlayersService.player.IsPlaying()
}

func (PlayersService *PlayersServiceStruct) Stop() bool {
	PlayersService.mut.RLock()
	if PlayersService.stopCtx == nil {
		return PlayersService.isPlaying
	}

	PlayersService.stopCtx()
	PlayersService.mut.RUnlock()

	PlayersService.updateIsPlaying()

	return PlayersService.isPlaying
}

func (PlayersService *PlayersServiceStruct) Pause() {
	PlayersService.mut.RLock()
	if PlayersService.player == nil {
		return
	}

	PlayersService.player.Pause()
	PlayersService.mut.RUnlock()

	PlayersService.updateIsPlaying()
}

func (PlayersService *PlayersServiceStruct) Resume() {
	PlayersService.mut.RLock()
	if PlayersService.player == nil {
		return
	}

	PlayersService.player.Play()
	PlayersService.mut.RUnlock()

	PlayersService.updateIsPlaying()
}

func (PlayersService *PlayersServiceStruct) PlayStationId(stationId string) (result model.PlayersData, err error) {
	PlayersService.Stop()

	err = PlayersService.updateDBAndPlayerService(stationId)
	result = model.PlayersData{}
	if err != nil {
		return result, err
	}

	err = PlayersService.playStream()
	if err != nil {
		return result, err
	}

	result, err = GetStationData()
	return result, err
}

func GetStationData() (result model.PlayersData, err error) {
	err = DB.Get(&result, "SELECT play, volume, stations.url AS url, stations.img AS img, stations.name AS name FROM state JOIN stations ON state.station_id = stations.id")
	if err != nil {
		log.Fatal("Impossible to get data in table state: ", err)
	}

	return result, nil
}

func (PlayersService *PlayersServiceStruct) SetVolume(volume float64) float64 {
	PlayersService.mut.Lock()
	if PlayersService.player == nil {
		return 0
	}

	PlayersService.player.SetVolume(volume)
	newVolume := PlayersService.player.Volume()
	PlayersService.volume = newVolume
	PlayersService.mut.Unlock()

	return newVolume
}
