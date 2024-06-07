package setting

import "time"

// 服务器配置
type ServerSettingS struct {
	RunMode      string
	HttpPort     string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

// Fisco目录配置
type FiscoSettingS struct {
	FiscoStartPath string
	FiscoPort      string
	FiscoStopPath  string
}

func (s *Setting) ReadSection(k string, v interface{}) error {
	err := s.vp.UnmarshalKey(k, v)
	if err != nil {
		return err
	}

	return nil
}
