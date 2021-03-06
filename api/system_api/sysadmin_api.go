package sysadmin

import (
	"context"
	_ "fmt"
	"perch/pkg/sysinfo"
	"perch/web/metric"
	"perch/web/model"
	"strconv"
	"strings"
	"time"

	"net/http"
)

func SysMemInfoHandler(w http.ResponseWriter, r *http.Request) {
	metric.ProcessMetricFunc(w, r, nil, func(ctx context.Context, bean interface{}, response *model.ResultReponse) error {
		var (
			sysMemInfo sysinfo.SysMemInformation
			err        error
		)

		response.Kind = "sysinfo memory"

		sysMemInfo, err = sysinfo.SysMemInfo()
		if err != nil {
			response.Code = http.StatusInternalServerError
			response.Message = err.Error()
			return err
		}

		response.Code = http.StatusOK
		response.Spec = sysMemInfo
		response.Message = " sys mem info"
		return nil
	})
}

func SysCpuInfoHandler(w http.ResponseWriter, req *http.Request) {
	metric.ProcessMetricFunc(w, req, nil, func(ctx context.Context, bean interface{}, response *model.ResultReponse) error {
		var (
			sysCpuInfo sysinfo.CpuAdvancedInfo
			logical    bool
			percpu     bool
			interval   time.Duration
			err        error
		)
		logicalStr := req.URL.Query().Get("logical")
		if logicalStr == "" {
			logical = true
		} else {
			logical, err = strconv.ParseBool(logicalStr)
			if err != nil {
				response.Code = http.StatusBadRequest
				response.Message = err.Error()
				return err
			}
		}

		percpuStr := req.URL.Query().Get("percpu")
		if percpuStr == "" {
			percpu = true
		} else {
			percpu, err = strconv.ParseBool(percpuStr)
			if err != nil {
				response.Code = http.StatusBadRequest
				response.Message = err.Error()
				return err
			}
		}

		intervalStr := req.URL.Query().Get("interval")
		if intervalStr == "" {
			interval = 1 * time.Second
		} else {
			interval, err = time.ParseDuration(intervalStr)
			if err != nil {
				response.Code = http.StatusBadRequest
				response.Message = err.Error()
				return err
			}
		}

		response.Kind = "sysinfo cpu"

		sysCpuInfo, err = sysinfo.SysAdvancedCpuInfo(logical, percpu, interval)
		if err != nil {
			response.Code = http.StatusInternalServerError
			response.Message = err.Error()
			return err
		}

		response.Code = http.StatusOK
		response.Spec = sysCpuInfo
		response.Total = 1
		response.Message = " sys mem info"
		return nil
	})
}

func SysHostInfoHandler(w http.ResponseWriter, r *http.Request) {
	metric.ProcessMetricFunc(w, r, nil, func(ctx context.Context, bean interface{}, response *model.ResultReponse) error {
		var (
			sysHostInfo sysinfo.HostAdvancedInfo
			err         error
		)

		response.Kind = "sysinfo memory"

		sysHostInfo, err = sysinfo.SysHostAdvancedInfo()
		if err != nil {
			response.Code = http.StatusInternalServerError
			response.Message = err.Error()
			return err
		}

		response.Code = http.StatusOK
		response.Spec = sysHostInfo
		response.Message = " sys host info"
		return nil
	})
}

func SysDockerInfoHandler(w http.ResponseWriter, r *http.Request) {
	metric.ProcessMetricFunc(w, r, nil, func(ctx context.Context, bean interface{}, response *model.ResultReponse) error {
		var (
			sysDockerInfo sysinfo.DockerAdvancedInfo
			err           error
		)

		response.Kind = "sysinfo docker"

		sysDockerInfo, err = sysinfo.SysAdvancedDockerInfo()
		if err != nil {
			response.Code = http.StatusInternalServerError
			response.Message = err.Error()
			return err
		}

		response.Code = http.StatusOK
		response.Spec = sysDockerInfo
		response.Total = 1
		response.Message = " sys docker info"
		return nil
	})
}

func SysDiskInfoHandler(w http.ResponseWriter, req *http.Request) {
	metric.ProcessMetricFunc(w, req, nil, func(ctx context.Context, bean interface{}, response *model.ResultReponse) error {
		var (
			diskSerialName string
			diskLableName  string
			partions       bool
			path           string
			iocounters     []string
			sysDiskInfo    sysinfo.DiskAdvacedInfo

			err error
		)

		diskSerialName = req.URL.Query().Get("diskSerialName")
		diskLableName = req.URL.Query().Get("diskLableName")
		partionsStr := req.URL.Query().Get("partions")
		if partionsStr == "" {
			partions = true
		} else {
			partions, err = strconv.ParseBool(partionsStr)
			if err != nil {
				response.Message = err.Error()
				response.Code = http.StatusBadRequest
				return err
			}

		}

		path = req.URL.Query().Get("path")
		iocountersStr := req.URL.Query().Get("iocounters")
		iocounters = strings.Split(iocountersStr, ",")
		response.Kind = "sysinfo disk"

		sysDiskInfo, err = sysinfo.SysAdvancedDiskInfo(diskSerialName, diskLableName, partions, path, iocounters...)
		if err != nil {
			response.Code = http.StatusInternalServerError
			response.Message = err.Error()
			return err
		}

		response.Total = 1
		response.Code = http.StatusOK
		response.Spec = sysDiskInfo
		response.Message = " sys disk info"
		return nil
	})
}

func SysNetInfoHandler(w http.ResponseWriter, req *http.Request) {
	metric.ProcessMetricFunc(w, req, nil, func(ctx context.Context, bean interface{}, response *model.ResultReponse) error {
		var (
			sysNetInfo sysinfo.NetAdvancedInfo
			percpu     bool
			err        error
		)

		response.Kind = "sysinfo net"
		percpuStr := req.URL.Query().Get("percpu")
		if percpuStr == "" {
			percpu = true
		} else {
			percpu, err = strconv.ParseBool(percpuStr)
			if err != nil {
				response.Code = http.StatusBadRequest
				response.Message = err.Error()
				return err
			}
		}
		sysNetInfo, err = sysinfo.SysAdvancedNetInfo(percpu)
		if err != nil {
			response.Code = http.StatusInternalServerError
			response.Message = err.Error()
			return err
		}

		response.Code = http.StatusOK
		response.Spec = sysNetInfo
		response.Message = " sys net info"
		return nil
	})
}

func SysProcessInfoHandler(w http.ResponseWriter, r *http.Request) {
	metric.ProcessMetricFunc(w, r, nil, func(ctx context.Context, bean interface{}, response *model.ResultReponse) error {
		var (
			sysProcessInfo sysinfo.ProcessAdvancedInfo
			err            error
		)

		response.Kind = "sysinfo process"

		sysProcessInfo, err = sysinfo.SysAdvancedProcessInfo()
		if err != nil {
			response.Code = http.StatusInternalServerError
			response.Message = err.Error()
			return err
		}

		response.Code = http.StatusOK
		response.Spec = sysProcessInfo
		response.Message = " sys process info"
		return nil
	})
}

func SysLoadInfoHandler(w http.ResponseWriter, r *http.Request) {
	metric.ProcessMetricFunc(w, r, nil, func(ctx context.Context, bean interface{}, response *model.ResultReponse) error {
		var (
			sysLoadInfo sysinfo.LoadAdvancedInfo
			err         error
		)

		response.Kind = "sysinfo memory"

		sysLoadInfo, err = sysinfo.SysAdvancedLoadInfo()
		if err != nil {
			response.Code = http.StatusInternalServerError
			response.Message = err.Error()
			return err
		}

		response.Code = http.StatusOK
		response.Spec = sysLoadInfo
		response.Message = " sys load info"
		return nil
	})
}
