package worker

// this file takes in a resource and returns a source (Volume)
// we might not need to model this way

import (
	"context"

	"github.com/concourse/concourse/atc/runtime"

	"code.cloudfoundry.org/lager"
	"github.com/concourse/concourse/atc"
	"github.com/concourse/concourse/atc/db"
	"github.com/concourse/concourse/atc/resource"
)

//go:generate counterfeiter . FetchSource

type FetchSource interface {
	//LockName() (string, error)
	Find() (Volume, bool, error)
	Create(context.Context) (GetResult, error)
}

//go:generate counterfeiter . FetchSourceFactory

type FetchSourceFactory interface {
	NewFetchSource(
		logger lager.Logger,
		worker Worker,
		resourceInstance resource.ResourceInstance,
		cache db.UsedResourceCache,
		resourceTypes atc.VersionedResourceTypes,
		containerSpec ContainerSpec,
		processSpec ProcessSpec,
		containerMetadata db.ContainerMetadata,
		imageFetchingDelegate ImageFetchingDelegate,
	) FetchSource
}

type fetchSourceFactory struct {
	resourceCacheFactory db.ResourceCacheFactory
	resourceFactory      resource.ResourceFactory
}

func NewFetchSourceFactory(
	resourceCacheFactory db.ResourceCacheFactory,
	resourceFactory resource.ResourceFactory,
) FetchSourceFactory {
	return &fetchSourceFactory{
		resourceCacheFactory: resourceCacheFactory,
		resourceFactory:      resourceFactory,
	}
}

func (r *fetchSourceFactory) NewFetchSource(
	logger lager.Logger,
	worker Worker,
	resourceInstance resource.ResourceInstance,
	cache db.UsedResourceCache,
	resourceTypes atc.VersionedResourceTypes,
	containerSpec ContainerSpec,
	processSpec ProcessSpec,
	containerMetadata db.ContainerMetadata,
	imageFetchingDelegate ImageFetchingDelegate,
) FetchSource {
	return &resourceInstanceFetchSource{
		logger:                 logger,
		worker:                 worker,
		resourceInstance:       resourceInstance,
		cache:                  cache,
		resourceTypes:          resourceTypes,
		containerSpec:          containerSpec,
		processSpec:            processSpec,
		containerMetadata:      containerMetadata,
		imageFetchingDelegate:  imageFetchingDelegate,
		dbResourceCacheFactory: r.resourceCacheFactory,
		resourceFactory:        r.resourceFactory,
	}
}

type resourceInstanceFetchSource struct {
	logger                 lager.Logger
	worker                 Worker
	resourceInstance       resource.ResourceInstance
	cache                  db.UsedResourceCache
	resourceTypes          atc.VersionedResourceTypes
	containerSpec          ContainerSpec
	processSpec            ProcessSpec
	containerMetadata      db.ContainerMetadata
	imageFetchingDelegate  ImageFetchingDelegate
	dbResourceCacheFactory db.ResourceCacheFactory
	resourceFactory        resource.ResourceFactory
}

//func (s *resourceInstanceFetchSource) LockName() (string, error) {
//	return s.resourceInstance.LockName(s.worker.Name())
//}

func findOn(logger lager.Logger, w Worker, cache db.UsedResourceCache) (volume Volume, found bool, err error) {
	return w.FindVolumeForResourceCache(
		logger,
		cache,
	)
}

func (s *resourceInstanceFetchSource) Find() (Volume, bool, error) {
	sLog := s.logger.Session("find")

	volume, found, err := findOn(s.logger, s.worker, s.cache)
	if err != nil {
		sLog.Error("failed-to-find-initialized-on", err)
		return nil, false, err
	}

	if !found {
		return nil, false, nil
	}

	//metadata, err := s.dbResourceCacheFactory.ResourceCacheMetadata(s.resourceInstance.ResourceCache())
	//if err != nil {
	//	sLog.Error("failed-to-get-resource-cache-metadata", err)
	//	return nil, false, err
	//}
	//
	//s.logger.Debug("found-initialized-versioned-source", lager.Data{"version": s.resourceInstance.Version(), "metadata": metadata.ToATCMetadata()})

	return volume,
		true, nil
}

// Create runs under the lock but we need to make sure volume does not exist
// yet before creating it under the lock
func (s *resourceInstanceFetchSource) Create(ctx context.Context) (GetResult, error) {
	sLog := s.logger.Session("create")
	result := GetResult{}

	volume, found, err := s.Find()
	if err != nil {
		return result, err
	}

	if found {
		result = GetResult{
			0,
			// todo: figure out what logically should be returned for VersionResult
			runtime.VersionResult{},
			runtime.GetArtifact{VolumeHandle: volume.Handle()},
			nil,
		}
		return result, nil
	}

	s.containerSpec.BindMounts = []BindMountSource{
		&CertsVolumeMount{Logger: s.logger},
	}

	container, err := s.worker.FindOrCreateContainer(
		ctx,
		s.logger,
		s.imageFetchingDelegate,
		s.resourceInstance.ContainerOwner(),
		s.containerMetadata,
		s.containerSpec,
		s.resourceTypes,
	)

	if err != nil {
		sLog.Error("failed-to-construct-resource", err)
		result = GetResult{
			1,
			// todo: figure out what logically should be returned for VersionResult
			runtime.VersionResult{},
			runtime.GetArtifact{VolumeHandle: volume.Handle()},
			err,
		}
		return result, err
	}

	mountPath := resource.ResourcesDir("get")
	for _, mount := range container.VolumeMounts() {
		if mount.MountPath == mountPath {
			volume = mount.Volume
			break
		}
	}

	vr := runtime.VersionResult{}
	events := make(chan runtime.Event)

	// todo: we want to decouple this resource from the container
	//res := s.resourceFactory.NewResourceForContainer(container)
	//versionedSource, err = res.Get(
	//	ctx,
	//	volume,
	//	runtime.IOConfig{
	//		Stdout: s.imageFetchingDelegate.Stdout(),
	//		Stderr: s.imageFetchingDelegate.Stderr(),
	//	},
	//	s.resourceInstance.Source(),
	//	s.resourceInstance.Params(),
	//	s.resourceInstance.Version(),
	//)
	//if err != nil {
	//	sLog.Error("failed-to-fetch-resource", err)
	//	return nil, err
	//}

	err = RunScript(
		ctx,
		container,
		s.processSpec.Path,
		s.processSpec.Args,
		runtime.GetRequest{
			Params: s.resourceInstance.Params(),
			Source: s.resourceInstance.Source(),
		},
		&vr,
		s.processSpec.StderrWriter,
		true,
		events,
	)

	if err != nil {
		sLog.Error("failed-to-fetch-resource", err)
		// TODO Is this compatible with previous behaviour of returning a nil when error type is NOT ErrResourceScriptFailed

		// if error returned from running the actual script
		if failErr, ok := err.(ErrResourceScriptFailed); ok {
			result = GetResult{failErr.ExitStatus, runtime.VersionResult{}, runtime.GetArtifact{}, failErr}
			return result, nil
		}
		return result, err
	}

	err = volume.SetPrivileged(false)
	if err != nil {
		sLog.Error("failed-to-set-volume-unprivileged", err)
		return result, err
	}

	// TODO this should happen get_step exec rather than here
	// seems like core logic

	//err = volume.InitializeResourceCache(s.resourceInstance.ResourceCache())
	//if err != nil {
	//	sLog.Error("failed-to-initialize-cache", err)
	//	return nil, err
	//}
	//
	//err = s.dbResourceCacheFactory.UpdateResourceCacheMetadata(s.resourceInstance.ResourceCache(), versionedSource.Metadata())
	//if err != nil {
	//	s.logger.Error("failed-to-update-resource-cache-metadata", err, lager.Data{"resource-cache": s.resourceInstance.ResourceCache()})
	//	return nil, err
	//}

	return GetResult{
		VersionResult: vr,
		GetArtifact: runtime.GetArtifact{
			VolumeHandle: volume.Handle(),
		},
	}, nil
}