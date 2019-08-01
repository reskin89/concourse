// Code generated by counterfeiter. DO NOT EDIT.
package dbfakes

import (
	"sync"

	"github.com/concourse/concourse/atc"
	"github.com/concourse/concourse/atc/db"
)

type FakeResourceCacheFactory struct {
	FindOrCreateResourceCacheStub        func(db.ResourceCacheUser, string, atc.Version, atc.Source, atc.Params, atc.VersionedResourceTypes) (db.UsedResourceCache, error)
	findOrCreateResourceCacheMutex       sync.RWMutex
	findOrCreateResourceCacheArgsForCall []struct {
		arg1 db.ResourceCacheUser
		arg2 string
		arg3 atc.Version
		arg4 atc.Source
		arg5 atc.Params
		arg6 atc.VersionedResourceTypes
	}
	findOrCreateResourceCacheReturns struct {
		result1 db.UsedResourceCache
		result2 error
	}
	findOrCreateResourceCacheReturnsOnCall map[int]struct {
		result1 db.UsedResourceCache
		result2 error
	}
	ResourceCacheMetadataStub        func(db.UsedResourceCache) (db.ResourceConfigMetadataFields, error)
	resourceCacheMetadataMutex       sync.RWMutex
	resourceCacheMetadataArgsForCall []struct {
		arg1 db.UsedResourceCache
	}
	resourceCacheMetadataReturns struct {
		result1 db.ResourceConfigMetadataFields
		result2 error
	}
	resourceCacheMetadataReturnsOnCall map[int]struct {
		result1 db.ResourceConfigMetadataFields
		result2 error
	}
	UpdateResourceCacheMetadataStub        func(db.UsedResourceCache, []atc.MetadataField) error
	updateResourceCacheMetadataMutex       sync.RWMutex
	updateResourceCacheMetadataArgsForCall []struct {
		arg1 db.UsedResourceCache
		arg2 []atc.MetadataField
	}
	updateResourceCacheMetadataReturns struct {
		result1 error
	}
	updateResourceCacheMetadataReturnsOnCall map[int]struct {
		result1 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeResourceCacheFactory) FindOrCreateResourceCache(arg1 db.ResourceCacheUser, arg2 string, arg3 atc.Version, arg4 atc.Source, arg5 atc.Params, arg6 atc.VersionedResourceTypes) (db.UsedResourceCache, error) {
	fake.findOrCreateResourceCacheMutex.Lock()
	ret, specificReturn := fake.findOrCreateResourceCacheReturnsOnCall[len(fake.findOrCreateResourceCacheArgsForCall)]
	fake.findOrCreateResourceCacheArgsForCall = append(fake.findOrCreateResourceCacheArgsForCall, struct {
		arg1 db.ResourceCacheUser
		arg2 string
		arg3 atc.Version
		arg4 atc.Source
		arg5 atc.Params
		arg6 atc.VersionedResourceTypes
	}{arg1, arg2, arg3, arg4, arg5, arg6})
	fake.recordInvocation("FindOrCreateResourceCache", []interface{}{arg1, arg2, arg3, arg4, arg5, arg6})
	fake.findOrCreateResourceCacheMutex.Unlock()
	if fake.FindOrCreateResourceCacheStub != nil {
		return fake.FindOrCreateResourceCacheStub(arg1, arg2, arg3, arg4, arg5, arg6)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.findOrCreateResourceCacheReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeResourceCacheFactory) FindOrCreateResourceCacheCallCount() int {
	fake.findOrCreateResourceCacheMutex.RLock()
	defer fake.findOrCreateResourceCacheMutex.RUnlock()
	return len(fake.findOrCreateResourceCacheArgsForCall)
}

func (fake *FakeResourceCacheFactory) FindOrCreateResourceCacheCalls(stub func(db.ResourceCacheUser, string, atc.Version, atc.Source, atc.Params, atc.VersionedResourceTypes) (db.UsedResourceCache, error)) {
	fake.findOrCreateResourceCacheMutex.Lock()
	defer fake.findOrCreateResourceCacheMutex.Unlock()
	fake.FindOrCreateResourceCacheStub = stub
}

func (fake *FakeResourceCacheFactory) FindOrCreateResourceCacheArgsForCall(i int) (db.ResourceCacheUser, string, atc.Version, atc.Source, atc.Params, atc.VersionedResourceTypes) {
	fake.findOrCreateResourceCacheMutex.RLock()
	defer fake.findOrCreateResourceCacheMutex.RUnlock()
	argsForCall := fake.findOrCreateResourceCacheArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3, argsForCall.arg4, argsForCall.arg5, argsForCall.arg6
}

func (fake *FakeResourceCacheFactory) FindOrCreateResourceCacheReturns(result1 db.UsedResourceCache, result2 error) {
	fake.findOrCreateResourceCacheMutex.Lock()
	defer fake.findOrCreateResourceCacheMutex.Unlock()
	fake.FindOrCreateResourceCacheStub = nil
	fake.findOrCreateResourceCacheReturns = struct {
		result1 db.UsedResourceCache
		result2 error
	}{result1, result2}
}

func (fake *FakeResourceCacheFactory) FindOrCreateResourceCacheReturnsOnCall(i int, result1 db.UsedResourceCache, result2 error) {
	fake.findOrCreateResourceCacheMutex.Lock()
	defer fake.findOrCreateResourceCacheMutex.Unlock()
	fake.FindOrCreateResourceCacheStub = nil
	if fake.findOrCreateResourceCacheReturnsOnCall == nil {
		fake.findOrCreateResourceCacheReturnsOnCall = make(map[int]struct {
			result1 db.UsedResourceCache
			result2 error
		})
	}
	fake.findOrCreateResourceCacheReturnsOnCall[i] = struct {
		result1 db.UsedResourceCache
		result2 error
	}{result1, result2}
}

func (fake *FakeResourceCacheFactory) ResourceCacheMetadata(arg1 db.UsedResourceCache) (db.ResourceConfigMetadataFields, error) {
	fake.resourceCacheMetadataMutex.Lock()
	ret, specificReturn := fake.resourceCacheMetadataReturnsOnCall[len(fake.resourceCacheMetadataArgsForCall)]
	fake.resourceCacheMetadataArgsForCall = append(fake.resourceCacheMetadataArgsForCall, struct {
		arg1 db.UsedResourceCache
	}{arg1})
	fake.recordInvocation("ResourceCacheMetadata", []interface{}{arg1})
	fake.resourceCacheMetadataMutex.Unlock()
	if fake.ResourceCacheMetadataStub != nil {
		return fake.ResourceCacheMetadataStub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.resourceCacheMetadataReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeResourceCacheFactory) ResourceCacheMetadataCallCount() int {
	fake.resourceCacheMetadataMutex.RLock()
	defer fake.resourceCacheMetadataMutex.RUnlock()
	return len(fake.resourceCacheMetadataArgsForCall)
}

func (fake *FakeResourceCacheFactory) ResourceCacheMetadataCalls(stub func(db.UsedResourceCache) (db.ResourceConfigMetadataFields, error)) {
	fake.resourceCacheMetadataMutex.Lock()
	defer fake.resourceCacheMetadataMutex.Unlock()
	fake.ResourceCacheMetadataStub = stub
}

func (fake *FakeResourceCacheFactory) ResourceCacheMetadataArgsForCall(i int) db.UsedResourceCache {
	fake.resourceCacheMetadataMutex.RLock()
	defer fake.resourceCacheMetadataMutex.RUnlock()
	argsForCall := fake.resourceCacheMetadataArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeResourceCacheFactory) ResourceCacheMetadataReturns(result1 db.ResourceConfigMetadataFields, result2 error) {
	fake.resourceCacheMetadataMutex.Lock()
	defer fake.resourceCacheMetadataMutex.Unlock()
	fake.ResourceCacheMetadataStub = nil
	fake.resourceCacheMetadataReturns = struct {
		result1 db.ResourceConfigMetadataFields
		result2 error
	}{result1, result2}
}

func (fake *FakeResourceCacheFactory) ResourceCacheMetadataReturnsOnCall(i int, result1 db.ResourceConfigMetadataFields, result2 error) {
	fake.resourceCacheMetadataMutex.Lock()
	defer fake.resourceCacheMetadataMutex.Unlock()
	fake.ResourceCacheMetadataStub = nil
	if fake.resourceCacheMetadataReturnsOnCall == nil {
		fake.resourceCacheMetadataReturnsOnCall = make(map[int]struct {
			result1 db.ResourceConfigMetadataFields
			result2 error
		})
	}
	fake.resourceCacheMetadataReturnsOnCall[i] = struct {
		result1 db.ResourceConfigMetadataFields
		result2 error
	}{result1, result2}
}

func (fake *FakeResourceCacheFactory) UpdateResourceCacheMetadata(arg1 db.UsedResourceCache, arg2 []atc.MetadataField) error {
	var arg2Copy []atc.MetadataField
	if arg2 != nil {
		arg2Copy = make([]atc.MetadataField, len(arg2))
		copy(arg2Copy, arg2)
	}
	fake.updateResourceCacheMetadataMutex.Lock()
	ret, specificReturn := fake.updateResourceCacheMetadataReturnsOnCall[len(fake.updateResourceCacheMetadataArgsForCall)]
	fake.updateResourceCacheMetadataArgsForCall = append(fake.updateResourceCacheMetadataArgsForCall, struct {
		arg1 db.UsedResourceCache
		arg2 []atc.MetadataField
	}{arg1, arg2Copy})
	fake.recordInvocation("UpdateResourceCacheMetadata", []interface{}{arg1, arg2Copy})
	fake.updateResourceCacheMetadataMutex.Unlock()
	if fake.UpdateResourceCacheMetadataStub != nil {
		return fake.UpdateResourceCacheMetadataStub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.updateResourceCacheMetadataReturns
	return fakeReturns.result1
}

func (fake *FakeResourceCacheFactory) UpdateResourceCacheMetadataCallCount() int {
	fake.updateResourceCacheMetadataMutex.RLock()
	defer fake.updateResourceCacheMetadataMutex.RUnlock()
	return len(fake.updateResourceCacheMetadataArgsForCall)
}

func (fake *FakeResourceCacheFactory) UpdateResourceCacheMetadataCalls(stub func(db.UsedResourceCache, []atc.MetadataField) error) {
	fake.updateResourceCacheMetadataMutex.Lock()
	defer fake.updateResourceCacheMetadataMutex.Unlock()
	fake.UpdateResourceCacheMetadataStub = stub
}

func (fake *FakeResourceCacheFactory) UpdateResourceCacheMetadataArgsForCall(i int) (db.UsedResourceCache, []atc.MetadataField) {
	fake.updateResourceCacheMetadataMutex.RLock()
	defer fake.updateResourceCacheMetadataMutex.RUnlock()
	argsForCall := fake.updateResourceCacheMetadataArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakeResourceCacheFactory) UpdateResourceCacheMetadataReturns(result1 error) {
	fake.updateResourceCacheMetadataMutex.Lock()
	defer fake.updateResourceCacheMetadataMutex.Unlock()
	fake.UpdateResourceCacheMetadataStub = nil
	fake.updateResourceCacheMetadataReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeResourceCacheFactory) UpdateResourceCacheMetadataReturnsOnCall(i int, result1 error) {
	fake.updateResourceCacheMetadataMutex.Lock()
	defer fake.updateResourceCacheMetadataMutex.Unlock()
	fake.UpdateResourceCacheMetadataStub = nil
	if fake.updateResourceCacheMetadataReturnsOnCall == nil {
		fake.updateResourceCacheMetadataReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.updateResourceCacheMetadataReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeResourceCacheFactory) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.findOrCreateResourceCacheMutex.RLock()
	defer fake.findOrCreateResourceCacheMutex.RUnlock()
	fake.resourceCacheMetadataMutex.RLock()
	defer fake.resourceCacheMetadataMutex.RUnlock()
	fake.updateResourceCacheMetadataMutex.RLock()
	defer fake.updateResourceCacheMetadataMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeResourceCacheFactory) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}

var _ db.ResourceCacheFactory = new(FakeResourceCacheFactory)
