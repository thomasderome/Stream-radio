package services

import (
	"context"
	"errors"
	"io"
	"log"
	"net/http"
	"os/exec"
	"radio_stream/model"
	"sync"
	"time"

	"github.com/ebitengine/oto/v3"
)

type PlayersServiceStruct struct {
	isPlaying bool
	volume    float64

	stationId   string
	stationUrl  string
	stationName string
	stationImg  string

	contextOto    *oto.Context
	player        *oto.Player
	stopCtxStream context.CancelFunc

	autoStopSteamCtx context.CancelFunc

	mut sync.RWMutex
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

	PlayersService.mut.Lock()
	defer PlayersService.mut.Unlock()
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
	player.SetVolume(PlayersService.volume)

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
	err = PlayersService.updateIsPlaying()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(context.Background())

	PlayersService.mut.Lock()
	PlayersService.stopCtxStream = cancel
	PlayersService.mut.Unlock()

	go PlayersService.keepAlive(ctx, resp, PlayersService.player, audio, cmd)

	return nil
}

func (PlayersService *PlayersServiceStruct) keepAlive(ctx context.Context, resp *http.Response, player *oto.Player, audio io.ReadCloser, cmd *exec.Cmd) {
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
	PlayersService.stopCtxStream = nil
}

func (PlayersService *PlayersServiceStruct) autoEndStreamLongPause(ctx context.Context) {
	timer := time.NewTimer(5 * time.Minute)
	defer timer.Stop()

	select {
	case <-timer.C:
		PlayersService.Stop()
		return
	case <-ctx.Done():
		return
	}
}

func (PlayersService *PlayersServiceStruct) updateIsPlaying() (err error) {
	PlayersService.mut.Lock()
	defer PlayersService.mut.Unlock()
	PlayersService.isPlaying = PlayersService.player.IsPlaying()

	_, err = DB.Exec("UPDATE state SET play = ? WHERE id = 1", PlayersService.isPlaying)
	if err != nil {
		return err
	}
	return nil
}

func (PlayersService *PlayersServiceStruct) Stop() (play bool, err error) {
	PlayersService.mut.RLock()

	if PlayersService.stopCtxStream == nil {
		PlayersService.mut.RUnlock()
		return PlayersService.isPlaying, nil
	}

	PlayersService.stopCtxStream()
	PlayersService.mut.RUnlock()

	err = PlayersService.updateIsPlaying()
	if err != nil {
		return false, err
	}
	return PlayersService.isPlaying, nil
}

func (PlayersService *PlayersServiceStruct) Pause() (err error) {
	PlayersService.mut.Lock()
	if PlayersService.player == nil {
		return errors.New("player is not created")
	}

	PlayersService.player.Pause()

	ctx, cancel := context.WithCancel(context.Background())
	go PlayersService.autoEndStreamLongPause(ctx)
	PlayersService.autoStopSteamCtx = cancel

	PlayersService.mut.Unlock()

	err = PlayersService.updateIsPlaying()
	if err != nil {
		return err
	}
	return nil
}

func (PlayersService *PlayersServiceStruct) Resume() (err error) {
	PlayersService.mut.Lock()

	if PlayersService.isPlaying {
		PlayersService.mut.Unlock()
		return errors.New("player is already playing")
	}

	if PlayersService.player == nil {
		PlayersService.mut.Unlock()
		err = PlayersService.playStream()
	} else if PlayersService.autoStopSteamCtx != nil {
		PlayersService.autoStopSteamCtx()
		PlayersService.autoStopSteamCtx = nil
		PlayersService.mut.Unlock()
	} else {
		PlayersService.mut.Unlock()
	}

	PlayersService.player.Play()

	err = PlayersService.updateIsPlaying()
	if err != nil {
		return err
	}

	return nil
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
	err = DB.Get(&result, "SELECT play, volume, COALESCE(stations.url, '') AS url, COALESCE(stations.img, '') AS img, COALESCE(stations.name, '') AS name FROM state LEFT JOIN stations ON state.station_id = stations.id")
	if err != nil {
		return result, err
	}

	return result, nil
}

func (PlayersService *PlayersServiceStruct) SetVolume(volume float64) (newVolume float64, err error) {
	PlayersService.mut.Lock()
	if PlayersService.player == nil && PlayersService.isPlaying {
		return 0, errors.New("player is not created")
	}
	if PlayersService.player != nil {
		PlayersService.player.SetVolume(volume)
		newVolume = PlayersService.player.Volume()
	} else {
		newVolume = volume
	}

	PlayersService.volume = newVolume
	PlayersService.mut.Unlock()

	_, err = DB.Exec("UPDATE state SET volume=? WHERE id=1", newVolume)
	if err != nil {
		return newVolume, err
	}

	return newVolume, nil
}
