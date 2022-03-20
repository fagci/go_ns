package svc

import (
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

type RTSP struct {
	Address string
	conn    net.Conn
	cseq    int
	ch      chan string
}

const PORT = "554"

const RTSP_HDR = "%s %s RTSP/1.0\r\n" +
	"CSeq: %d\r\n" +
	"User-Agent: LibVLC/3.0.0\r\n" +
	"Accept: application/sdp\r\n\r\n"

func (r *RTSP) Request(req string) (int, error) {
	if _, e := r.conn.Write([]byte(req)); e != nil {
		return 0, e
	}

	m := make([]byte, 1024)
	if _, e := r.conn.Read(m); e != nil {
		return 0, e
	}

	f := strings.Fields(string(m))
	if len(f) > 2 && strings.HasPrefix(f[0], "RTSP") {
		return strconv.Atoi(f[1])
	}

	return 0, errors.New("Bad response")
}

func (r *RTSP) Query(path string) string {
	var method string

	if path == "*" {
		method = "OPTIONS"
	} else {
		method = "DESCRIBE"
	}

	return fmt.Sprintf(RTSP_HDR, method, path, r.cseq)
}

func (r *RTSP) check(paths *[]string) {
	defer close(r.ch)
	d := net.Dialer{Timeout: time.Second * 2}
	var err error

	r.conn, err = d.Dial("tcp", r.Address)

	if err != nil {
		return
	}

	defer r.conn.Close()

	var code int
	code, err = r.Request(r.Query("*"))
	if err != nil {
		return
	}

	code, err = r.Request(r.Query("/"))
	if err != nil || code == 401 {
		return
	}

	if code == 200 {
		r.ch <- fmt.Sprintf("rtsp://%s/", r.Address)
		return
	}

	for _, path := range *paths {
		code, err = r.Request(r.Query(path))
		if err != nil || code == 401 {
			return
		}
		if code == 200 {
			r.ch <- fmt.Sprintf("rtsp://%s%s", r.Address, path)
			return
		}
	}
}

func (r *RTSP) CheckPaths(paths *[]string) <-chan string {
	go r.check(paths)
	return r.ch
}

func NewRTSP(address string) *RTSP {
	return &RTSP{
		Address: address,
		ch:      make(chan string),
	}
}

var RTSP_PATHS = []string{
	"/1",
	"/0/1:1/main",
	"/live/h264",
	"/live",
	"/h264/ch1/sub/av_stream",
	"/stream1",
	"/live.sdp",
	"/image.mpg",
	"/axis-media/media.amp",
	"/1/stream1",
	"/ch01.264",
	"/live1.sdp",
	"/stream.sdp",
	"/0/usrnm:pwd/main",
	"/0/video1",
	"/1.AMP",
	"/1/h264major",
	"/1080p",
	"/11",
	"/12",
	"/125",
	"/1440p",
	"/480p",
	"/4K",
	"/666",
	"/720p",
	"/AVStream1_1",
	"/CH001.sdp",
	"/GetData.cgi",
	"/HD",
	"/HighResolutionVideo",
	"/LowResolutionVideo",
	"/MediaInput/h264",
	"/MediaInput/mpeg4",
	"/ONVIF/MediaInput",
	"/ONVIF/MediaInput?profile=4_def_profile6",
	"/StdCh1",
	"/Streaming/Channels/1",
	"/Streaming/Unicast/channels/101",
	"/StreamingSetting?version=1.0&action=getRTSPStream&ChannelID=1&ChannelName=Channel1",
	"/VideoInput/1/h264/1",
	"/VideoInput/1/mpeg4/1",
	"/access_code",
	"/access_name_for_stream_1_to_5",
	"/api/mjpegvideo.cgi",
	"/av0_0",
	"/av2",
	"/avc",
	"/avn=2",
	"/axis-media/media.amp?camera=1",
	"/axis-media/media.amp?videocodec=h264",
	"/cam",
	"/cam/realmonitor",
	"/cam/realmonitor?channel=0&subtype=0",
	"/cam/realmonitor?channel=1&subtype=0",
	"/cam/realmonitor?channel=1&subtype=1",
	"/cam/realmonitor?channel=1&subtype=1&unicast=true&proto=Onvif",
	"/cam0",
	"/cam0_0",
	"/cam0_1",
	"/cam1",
	"/cam1/h264",
	"/cam1/h264/multicast",
	"/cam1/mjpeg",
	"/cam1/mpeg4",
	"/cam1/onvif-h264",
	"/camera.stm",
	"/ch0",
	"/ch00/0",
	"/ch001.sdp",
	"/ch01.264?",
	"/ch01.264?ptype=tcp",
	"/ch0_0.h264",
	"/ch0_unicast_firststream",
	"/ch0_unicast_secondstream",
	"/ch1-s1",
	"/ch1/0",
	"/ch1_0",
	"/ch2/0",
	"/ch2_0",
	"/ch3/0",
	"/ch3_0",
	"/ch4/0",
	"/ch4_0",
	"/channel1",
	"/gnz_media/main",
	"/h264",
	"/h264.sdp",
	"/h264/media.amp",
	"/h264Preview_01_main",
	"/h264Preview_01_sub",
	"/h264_stream",
	"/h264_vga.sdp",
	"/img/media.sav",
	"/img/media.sav?channel=1",
	"/img/video.asf",
	"/img/video.sav",
	"/ioImage/1",
	"/ipcam.sdp",
	"/ipcam_h264.sdp",
	"/ipcam_mjpeg.sdp",
	"/live/av0",
	"/live/ch0",
	"/live/ch00_0",
	"/live/ch01_0",
	"/live/main",
	"/live/main0",
	"/live/mpeg4",
	"/live3.sdp",
	"/live_mpeg4.sdp",
	"/live_st1",
	"/livestream",
	"/main",
	"/media",
	"/media.amp",
	"/media.amp?streamprofile=Profile1",
	"/media/media.amp",
	"/media/video1",
	"/medias2",
	"/mjpeg/media.smp",
	"/mp4",
	"/mpeg/media.amp",
	"/mpeg4",
	"/mpeg4/1/media.amp",
	"/mpeg4/media.amp",
	"/mpeg4/media.smp",
	"/mpeg4unicast",
	"/mpg4/rtsp.amp",
	"/multicaststream",
	"/now.mp4",
	"/nph-h264.cgi",
	"/nphMpeg4/g726-640x",
	"/nphMpeg4/g726-640x48",
	"/nphMpeg4/g726-640x480",
	"/nphMpeg4/nil-320x240",
	"/onvif-media/media.amp",
	"/onvif1",
	"/play1.sdp",
	"/play2.sdp",
	"/profile2/media.smp",
	"/profile5/media.smp",
	"/rtpvideo1.sdp",
	"/rtsp_live0",
	"/rtsp_live1",
	"/rtsp_live2",
	"/rtsp_tunnel",
	"/rtsph264",
	"/rtsph2641080p",
	"/snap.jpg",
	"/stream",
	"/stream/0",
	"/stream/1",
	"/stream/live.sdp",
	"/streaming/channels/0",
	"/streaming/channels/1",
	"/streaming/channels/101",
	"/tcp/av0_0",
	"/test",
	"/tmpfs/auto.jpg",
	"/trackID=1",
	"/ucast/11",
	"/udp/av0_0",
	"/udp/unicast/aiphone_H264",
	"/udpstream",
	"/user.pin.mp2",
	"/user_defined",
	"/v2",
	"/video",
	"/video.3gp",
	"/video.h264",
	"/video.mjpg",
	"/video.mp4",
	"/video.pro1",
	"/video.pro2",
	"/video.pro3",
	"/video0",
	"/video0.sdp",
	"/video1",
	"/video1+audio1",
	"/video1.sdp",
	"/videoMain",
	"/videoinput_1/h264_1/media.stm",
	"/videostream.asf",
	"/vis",
	"/wfov",
}
