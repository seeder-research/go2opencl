package go2opencl

/*
#include "./opencl.h"


static cl_device_partition_property * partitionDeviceEqually(unsigned int n) {
	cl_device_partition_property *properties = malloc(n * sizeof(cl_device_partition_property));
	properties[0] = CL_DEVICE_PARTITION_EQUALLY;
	properties[1] = (cl_device_partition_property)(n);
	properties[2] = CL_DEVICE_PARTITION_BY_COUNTS_LIST_END;
	return properties;
}

static cl_device_partition_property * partitionDeviceByCounts(unsigned int *n, unsigned int num_counts) {
	cl_device_partition_property *properties = malloc((num_counts + 3) * sizeof(cl_device_partition_property));
	properties[0] = CL_DEVICE_PARTITION_BY_COUNTS;
	int idx;
	for (idx = 0; idx < num_counts; idx++) {
		properties[idx+1] = (cl_device_partition_property)(*(n+idx));
	}
	properties[num_counts+1] = CL_DEVICE_PARTITION_BY_COUNTS_LIST_END;
	properties[num_counts+2] = 0;
	return properties;
}

static cl_device_partition_property * partitionDeviceByNuma() {
        cl_device_partition_property *properties = malloc(3 * sizeof(cl_device_partition_property));
        properties[0] = CL_DEVICE_PARTITION_EQUALLY;
        properties[1] = (cl_device_partition_property)(CL_DEVICE_AFFINITY_DOMAIN_NUMA);
        properties[2] = CL_DEVICE_PARTITION_BY_COUNTS_LIST_END;
        return properties;
}

static cl_device_partition_property * partitionDeviceByL4Cache() {
        cl_device_partition_property *properties = malloc(3 * sizeof(cl_device_partition_property));
        properties[0] = CL_DEVICE_PARTITION_EQUALLY;
        properties[1] = (cl_device_partition_property)(CL_DEVICE_AFFINITY_DOMAIN_L4_CACHE);
        properties[2] = CL_DEVICE_PARTITION_BY_COUNTS_LIST_END;
        return properties;
}

static cl_device_partition_property * partitionDeviceByL3Cache() {
        cl_device_partition_property *properties = malloc(3 * sizeof(cl_device_partition_property));
        properties[0] = CL_DEVICE_PARTITION_EQUALLY;
        properties[1] = (cl_device_partition_property)(CL_DEVICE_AFFINITY_DOMAIN_L3_CACHE);
        properties[2] = CL_DEVICE_PARTITION_BY_COUNTS_LIST_END;
        return properties;
}

static cl_device_partition_property * partitionDeviceByL2Cache() {
        cl_device_partition_property *properties = malloc(3 * sizeof(cl_device_partition_property));
        properties[0] = CL_DEVICE_PARTITION_EQUALLY;
        properties[1] = (cl_device_partition_property)(CL_DEVICE_AFFINITY_DOMAIN_L2_CACHE);
        properties[2] = CL_DEVICE_PARTITION_BY_COUNTS_LIST_END;
        return properties;
}

static cl_device_partition_property * partitionDeviceByL1Cache() {
        cl_device_partition_property *properties = malloc(3 * sizeof(cl_device_partition_property));
        properties[0] = CL_DEVICE_PARTITION_EQUALLY;
        properties[1] = (cl_device_partition_property)(CL_DEVICE_AFFINITY_DOMAIN_L1_CACHE);
        properties[2] = CL_DEVICE_PARTITION_BY_COUNTS_LIST_END;
        return properties;
}

static cl_device_partition_property * partitionDeviceByNextPartitionable() {
        cl_device_partition_property *properties = malloc(3 * sizeof(cl_device_partition_property));
        properties[0] = CL_DEVICE_PARTITION_EQUALLY;
        properties[1] = (cl_device_partition_property)(CL_DEVICE_AFFINITY_DOMAIN_NEXT_PARTITIONABLE);
        properties[2] = CL_DEVICE_PARTITION_BY_COUNTS_LIST_END;
        return properties;
}

static cl_int CLGetDeviceInfoParamSize(cl_device_id device, cl_device_info param_name, size_t* param_value_size_ret) {
        return clGetDeviceInfo(device, param_name, NULL, NULL, param_value_size_ret);
}

static cl_int CLGetDeviceInfoParamUnsafe(cl_device_id device, cl_device_info param_name, size_t param_value_size, void *param_value) {
        return clGetDeviceInfo(device, param_name, param_value_size, param_value, NULL);
}

*/
import "C"

import (
	"strings"
	"unsafe"
)

// Unsupported device info queries:
// CL_DEVICE_PARTITION_PROPERTIES
// CL_DEVICE_PARTITION_TYPE

//////////////// Constants ////////////////
const maxDeviceCount = 64

//////////////// Basic Types ////////////////
type DeviceType uint

const (
	DeviceTypeCPU         DeviceType = C.CL_DEVICE_TYPE_CPU
	DeviceTypeGPU         DeviceType = C.CL_DEVICE_TYPE_GPU
	DeviceTypeAccelerator DeviceType = C.CL_DEVICE_TYPE_ACCELERATOR
	DeviceTypeCustom      DeviceType = C.CL_DEVICE_TYPE_CUSTOM
	DeviceTypeDefault     DeviceType = C.CL_DEVICE_TYPE_DEFAULT
	DeviceTypeAll         DeviceType = C.CL_DEVICE_TYPE_ALL
)

type FPConfig int

const (
	FPConfigDenorm         FPConfig = C.CL_FP_DENORM           // denorms are supported
	FPConfigInfNaN         FPConfig = C.CL_FP_INF_NAN          // INF and NaNs are supported
	FPConfigRoundToNearest FPConfig = C.CL_FP_ROUND_TO_NEAREST // round to nearest even rounding mode supported
	FPConfigRoundToZero    FPConfig = C.CL_FP_ROUND_TO_ZERO    // round to zero rounding mode supported
	FPConfigRoundToInf     FPConfig = C.CL_FP_ROUND_TO_INF     // round to positive and negative infinity rounding modes supported
	FPConfigFMA            FPConfig = C.CL_FP_FMA              // IEEE754-2008 fused multiply-add is supported
)

type DeviceAffinityDomain uint

const (
	DeviceAffinityDomainNuma              DeviceAffinityDomain = C.CL_DEVICE_AFFINITY_DOMAIN_NUMA
	DeviceAffinityDomainL4Cache           DeviceAffinityDomain = C.CL_DEVICE_AFFINITY_DOMAIN_L4_CACHE
	DeviceAffinityDomainL3Cache           DeviceAffinityDomain = C.CL_DEVICE_AFFINITY_DOMAIN_L3_CACHE
	DeviceAffinityDomainL2Cache           DeviceAffinityDomain = C.CL_DEVICE_AFFINITY_DOMAIN_L2_CACHE
	DeviceAffinityDomainL1Cache           DeviceAffinityDomain = C.CL_DEVICE_AFFINITY_DOMAIN_L1_CACHE
	DeviceAffinityDomainNextPartitionable DeviceAffinityDomain = C.CL_DEVICE_AFFINITY_DOMAIN_NEXT_PARTITIONABLE
)

var fpConfigNameMap = map[FPConfig]string{
	FPConfigDenorm:         "Denorm",
	FPConfigInfNaN:         "InfNaN",
	FPConfigRoundToNearest: "RoundToNearest",
	FPConfigRoundToZero:    "RoundToZero",
	FPConfigRoundToInf:     "RoundToInf",
	FPConfigFMA:            "FMA",
}

func (c FPConfig) String() string {
	var parts []string
	for bit, name := range fpConfigNameMap {
		if c&bit != 0 {
			parts = append(parts, name)
		}
	}
	if parts == nil {
		return ""
	}
	return strings.Join(parts, "|")
}

func (dt DeviceType) String() string {
	var parts []string
	if dt&DeviceTypeCPU != 0 {
		parts = append(parts, "CPU")
	}
	if dt&DeviceTypeGPU != 0 {
		parts = append(parts, "GPU")
	}
	if dt&DeviceTypeAccelerator != 0 {
		parts = append(parts, "Accelerator")
	}
	if dt&DeviceTypeDefault != 0 {
		parts = append(parts, "Default")
	}
	if parts == nil {
		parts = append(parts, "None")
	}
	return strings.Join(parts, "|")
}

//////////////// Abstract Types ////////////////
type Device struct {
	id C.cl_device_id
}

//////////////// Golang Types ////////////////
type CLDevice C.cl_device_id

//////////////// Basic Functions ////////////////
func buildDeviceIdList(devices []*Device) []C.cl_device_id {
	deviceIds := make([]C.cl_device_id, len(devices))
	for i, d := range devices {
		deviceIds[i] = d.id
	}
	return deviceIds
}

// Obtain the list of devices available on a platform. 'platform' refers
// to the platform returned by GetPlatforms or can be nil. If platform
// is nil, the behavior is implementation-defined.
func GetDevices(platform *Platform, deviceType DeviceType) ([]*Device, error) {
	var deviceIds [maxDeviceCount]C.cl_device_id
	var numDevices C.cl_uint
	var platformId C.cl_platform_id
	if platform != nil {
		platformId = platform.id
	}
	if err := C.clGetDeviceIDs(platformId, C.cl_device_type(deviceType), C.cl_uint(maxDeviceCount), &deviceIds[0], &numDevices); err != C.CL_SUCCESS {
		return nil, toError(err)
	}
	if numDevices > maxDeviceCount {
		numDevices = maxDeviceCount
	}
	devices := make([]*Device, numDevices)
	for i := 0; i < int(numDevices); i++ {
		devices[i] = &Device{id: deviceIds[i]}
	}
	return devices, nil
}

//////////////// Abstract Functions ////////////////
func (d *Device) nullableId() C.cl_device_id {
	if d == nil {
		return nil
	}
	return d.id
}

func (d *Device) GetInfoString(param C.cl_device_info, panicOnError bool) (string, error) {
	var strN C.size_t
	if err := C.CLGetDeviceInfoParamSize(d.nullableId(), param, &strN); err != C.CL_SUCCESS {
		if panicOnError {
			panic("Should never fail getting size of parameter")
		}
		return "", toError(err)
	}
	strC := (*C.char)(C.calloc(strN, 1))
	defer C.free(unsafe.Pointer(strC))
	if err := C.CLGetDeviceInfoParamUnsafe(d.nullableId(), param, strN, unsafe.Pointer(strC)); err != C.CL_SUCCESS {
		if panicOnError {
			panic("Should never fail getting device info")
		}
		return "", toError(err)
	}

	// OpenCL strings are NUL-terminated, and the terminator is included in strN
	// Go strings aren't NUL-terminated, so subtract 1 from the length
	retString := C.GoStringN(strC, C.int(strN-1))
	return retString, nil
}

func (d *Device) getInfoUint(param C.cl_device_info, panicOnError bool) (uint, error) {
	var val C.cl_uint
	if err := C.clGetDeviceInfo(d.nullableId(), param, C.size_t(unsafe.Sizeof(val)), unsafe.Pointer(&val), nil); err != C.CL_SUCCESS {
		if panicOnError {
			panic("Should never fail")
		}
		return 0, toError(err)
	}
	return uint(val), nil
}

func (d *Device) getInfoSize(param C.cl_device_info, panicOnError bool) (int, error) {
	var val C.size_t
	if err := C.clGetDeviceInfo(d.nullableId(), param, C.size_t(unsafe.Sizeof(val)), unsafe.Pointer(&val), nil); err != C.CL_SUCCESS {
		if panicOnError {
			panic("Should never fail")
		}
		return 0, toError(err)
	}
	return int(val), nil
}

func (d *Device) getInfoUlong(param C.cl_device_info, panicOnError bool) (int64, error) {
	var val C.cl_ulong
	if err := C.clGetDeviceInfo(d.nullableId(), param, C.size_t(unsafe.Sizeof(val)), unsafe.Pointer(&val), nil); err != C.CL_SUCCESS {
		if panicOnError {
			panic("Should never fail")
		}
		return 0, toError(err)
	}
	return int64(val), nil
}

func (d *Device) getInfoBool(param C.cl_device_info, panicOnError bool) (bool, error) {
	var val C.cl_bool
	if err := C.clGetDeviceInfo(d.nullableId(), param, C.size_t(unsafe.Sizeof(val)), unsafe.Pointer(&val), nil); err != C.CL_SUCCESS {
		if panicOnError {
			panic("Should never fail")
		}
		return false, toError(err)
	}
	return val == C.CL_TRUE, nil
}

func (d *Device) Name() string {
	str, _ := d.GetInfoString(C.CL_DEVICE_NAME, true)
	return str
}

func (d *Device) Platform() *Platform {
	var devicePlatform C.cl_platform_id
	if err := C.clGetDeviceInfo(d.nullableId(), C.CL_DEVICE_PLATFORM, C.size_t(unsafe.Sizeof(devicePlatform)), unsafe.Pointer(&devicePlatform), nil); err != C.CL_SUCCESS {
		panic("Failed to get device platform")
	}
	return &Platform{id: devicePlatform}
}

func (d *Device) Vendor() string {
	str, _ := d.GetInfoString(C.CL_DEVICE_VENDOR, true)
	return str
}

func (d *Device) VendorId() int {
	val, _ := d.getInfoUint(C.CL_DEVICE_VENDOR_ID, true)
	return int(val)
}

func (d *Device) Extensions() string {
	str, _ := d.GetInfoString(C.CL_DEVICE_EXTENSIONS, true)
	return str
}

func (d *Device) Profile() string {
	str, _ := d.GetInfoString(C.CL_DEVICE_PROFILE, true)
	return str
}

func (d *Device) Version() string {
	str, _ := d.GetInfoString(C.CL_DEVICE_VERSION, true)
	return str
}

func (d *Device) DriverVersion() string {
	str, _ := d.GetInfoString(C.CL_DRIVER_VERSION, true)
	return str
}

func (d *Device) OpenCLCVersion() string {
	str, _ := d.GetInfoString(C.CL_DEVICE_OPENCL_C_VERSION, true)
	return str
}

// Built-in kernels supported by the device
func (d *Device) BuiltInKernels() string {
	str, _ := d.GetInfoString(C.CL_DEVICE_BUILT_IN_KERNELS, true)
	return str
}

// Device reference count
func (d *Device) ReferenceCount() int {
	val, _ := d.getInfoUint(C.CL_DEVICE_REFERENCE_COUNT, true)
	return int(val)
}

// The default compute device address space size specified as an
// unsigned integer value in bits. Currently supported values are 32 or 64 bits.
func (d *Device) AddressBits() int {
	val, _ := d.getInfoUint(C.CL_DEVICE_ADDRESS_BITS, true)
	return int(val)
}

// Size of global memory cache line in bytes.
func (d *Device) GlobalMemCachelineSize() int {
	val, _ := d.getInfoUint(C.CL_DEVICE_GLOBAL_MEM_CACHELINE_SIZE, true)
	return int(val)
}

// Maximum configured clock frequency of the device in MHz.
func (d *Device) MaxClockFrequency() int {
	val, _ := d.getInfoUint(C.CL_DEVICE_MAX_CLOCK_FREQUENCY, true)
	return int(val)
}

// The number of parallel compute units on the OpenCL device.
// A work-group executes on a single compute unit. The minimum value is 1.
func (d *Device) MaxComputeUnits() int {
	val, _ := d.getInfoUint(C.CL_DEVICE_MAX_COMPUTE_UNITS, true)
	return int(val)
}

// Max number of arguments declared with the __constant qualifier in a kernel.
// The minimum value is 8 for devices that are not of type CL_DEVICE_TYPE_CUSTOM.
func (d *Device) MaxConstantArgs() int {
	val, _ := d.getInfoUint(C.CL_DEVICE_MAX_CONSTANT_ARGS, true)
	return int(val)
}

// Max number of images in a 1D or 2D image array. The minimum value is 2048
// if CL_DEVICE_IMAGE_SUPPORT is CL_TRUE
func (d *Device) MaxImageArraySize() int {
	val, _ := d.getInfoUint(C.CL_DEVICE_IMAGE_MAX_ARRAY_SIZE, true)
	return int(val)
}

// Max number of pixels for a 1D image created from a buffer object.
// The minimum value is 65536 if CL_DEVICE_IMAGE_SUPPORT is CL_TRUE.
func (d *Device) MaxImageBufferSize() int {
	val, _ := d.getInfoUint(C.CL_DEVICE_IMAGE_MAX_BUFFER_SIZE, true)
	return int(val)
}

// Max number of simultaneous image objects that can be read by a kernel.
// The minimum value is 128 if CL_DEVICE_IMAGE_SUPPORT is CL_TRUE.
func (d *Device) MaxReadImageArgs() int {
	val, _ := d.getInfoUint(C.CL_DEVICE_MAX_READ_IMAGE_ARGS, true)
	return int(val)
}

// Maximum number of samplers that can be used in a kernel. The minimum
// value is 16 if CL_DEVICE_IMAGE_SUPPORT is CL_TRUE. (Also see sampler_t.)
func (d *Device) MaxSamplers() int {
	val, _ := d.getInfoUint(C.CL_DEVICE_MAX_SAMPLERS, true)
	return int(val)
}

// Maximum dimensions that specify the global and local work-item IDs used
// by the data parallel execution model. (Refer to clEnqueueNDRangeKernel).
// The minimum value is 3 for devices that are not of type CL_DEVICE_TYPE_CUSTOM.
func (d *Device) MaxWorkItemDimensions() int {
	val, _ := d.getInfoUint(C.CL_DEVICE_MAX_WORK_ITEM_DIMENSIONS, true)
	return int(val)
}

// Max number of simultaneous image objects that can be written to by a
// kernel. The minimum value is 8 if CL_DEVICE_IMAGE_SUPPORT is CL_TRUE.
func (d *Device) MaxWriteImageArgs() int {
	val, _ := d.getInfoUint(C.CL_DEVICE_MAX_WRITE_IMAGE_ARGS, true)
	return int(val)
}

// The minimum value is the size (in bits) of the largest OpenCL built-in
// data type supported by the device (long16 in FULL profile, long16 or
// int16 in EMBEDDED profile) for devices that are not of type CL_DEVICE_TYPE_CUSTOM.
func (d *Device) MemBaseAddrAlign() int {
	val, _ := d.getInfoUint(C.CL_DEVICE_MEM_BASE_ADDR_ALIGN, true)
	return int(val)
}

// Min of bytes for the smallest alignment that can be used for any data type
func (d *Device) MinDataTypeAlignSize() int {
	val, _ := d.getInfoUint(C.CL_DEVICE_MIN_DATA_TYPE_ALIGN_SIZE, true)
	return int(val)
}

// Max height of 2D image in pixels. The minimum value is 8192
// if CL_DEVICE_IMAGE_SUPPORT is CL_TRUE.
func (d *Device) Image2DMaxHeight() int {
	val, _ := d.getInfoSize(C.CL_DEVICE_IMAGE2D_MAX_HEIGHT, true)
	return int(val)
}

// Max width of 2D image or 1D image not created from a buffer object in
// pixels. The minimum value is 8192 if CL_DEVICE_IMAGE_SUPPORT is CL_TRUE.
func (d *Device) Image2DMaxWidth() int {
	val, _ := d.getInfoSize(C.CL_DEVICE_IMAGE2D_MAX_WIDTH, true)
	return int(val)
}

// Max depth of 3D image in pixels. The minimum value is 2048 if CL_DEVICE_IMAGE_SUPPORT is CL_TRUE.
func (d *Device) Image3DMaxDepth() int {
	val, _ := d.getInfoSize(C.CL_DEVICE_IMAGE3D_MAX_DEPTH, true)
	return int(val)
}

// Max height of 3D image in pixels. The minimum value is 2048 if CL_DEVICE_IMAGE_SUPPORT is CL_TRUE.
func (d *Device) Image3DMaxHeight() int {
	val, _ := d.getInfoSize(C.CL_DEVICE_IMAGE3D_MAX_HEIGHT, true)
	return int(val)
}

// Max width of 3D image in pixels. The minimum value is 2048 if CL_DEVICE_IMAGE_SUPPORT is CL_TRUE.
func (d *Device) Image3DMaxWidth() int {
	val, _ := d.getInfoSize(C.CL_DEVICE_IMAGE3D_MAX_WIDTH, true)
	return int(val)
}

// Max size in bytes of the arguments that can be passed to a kernel. The
// minimum value is 1024 for devices that are not of type CL_DEVICE_TYPE_CUSTOM.
// For this minimum value, only a maximum of 128 arguments can be passed to a kernel.
func (d *Device) MaxParameterSize() int {
	val, _ := d.getInfoSize(C.CL_DEVICE_MAX_PARAMETER_SIZE, true)
	return int(val)
}

// Maximum number of work-items in a work-group executing a kernel on a
// single compute unit, using the data parallel execution model. (Refer
// to clEnqueueNDRangeKernel). The minimum value is 1.
func (d *Device) MaxWorkGroupSize() int {
	val, _ := d.getInfoSize(C.CL_DEVICE_MAX_WORK_GROUP_SIZE, true)
	return int(val)
}

// Describes the resolution of device timer. This is measured in nanoseconds.
func (d *Device) ProfilingTimerResolution() int {
	val, _ := d.getInfoSize(C.CL_DEVICE_PROFILING_TIMER_RESOLUTION, true)
	return int(val)
}

// Maximum size of the internal buffer that holds the output of printf calls from a
// kernel. The minimum value for the FULL profile is 1 MB.
func (d *Device) PrintfBufferSize() int {
	val, _ := d.getInfoSize(C.CL_DEVICE_PRINTF_BUFFER_SIZE, true)
	return int(val)
}

// Size of local memory arena in bytes. The minimum value is 32 KB for
// devices that are not of type CL_DEVICE_TYPE_CUSTOM.
func (d *Device) LocalMemSize() int64 {
	val, _ := d.getInfoUlong(C.CL_DEVICE_LOCAL_MEM_SIZE, true)
	return val
}

// Max size in bytes of a constant buffer allocation. The minimum value is
// 64 KB for devices that are not of type CL_DEVICE_TYPE_CUSTOM.
func (d *Device) MaxConstantBufferSize() int64 {
	val, _ := d.getInfoUlong(C.CL_DEVICE_MAX_CONSTANT_BUFFER_SIZE, true)
	return val
}

// Max size of memory object allocation in bytes. The minimum value is max
// (1/4th of CL_DEVICE_GLOBAL_MEM_SIZE, 128*1024*1024) for devices that are
// not of type CL_DEVICE_TYPE_CUSTOM.
func (d *Device) MaxMemAllocSize() int64 {
	val, _ := d.getInfoUlong(C.CL_DEVICE_MAX_MEM_ALLOC_SIZE, true)
	return val
}

// Size of global device memory in bytes.
func (d *Device) GlobalMemSize() int64 {
	val, _ := d.getInfoUlong(C.CL_DEVICE_GLOBAL_MEM_SIZE, true)
	return val
}

// Size of global device cache in bytes.
func (d *Device) GlobalMemCacheSize() int64 {
	val, _ := d.getInfoUlong(C.CL_DEVICE_GLOBAL_MEM_CACHE_SIZE, true)
	return val
}

func (d *Device) Available() bool {
	val, _ := d.getInfoBool(C.CL_DEVICE_AVAILABLE, true)
	return val
}

func (d *Device) HostUnifiedMemory() bool {
	val, _ := d.getInfoBool(C.CL_DEVICE_HOST_UNIFIED_MEMORY, true)
	return val
}

func (d *Device) CompilerAvailable() bool {
	val, _ := d.getInfoBool(C.CL_DEVICE_COMPILER_AVAILABLE, true)
	return val
}

func (d *Device) LinkerAvailable() bool {
	val, _ := d.getInfoBool(C.CL_DEVICE_LINKER_AVAILABLE, true)
	return val
}

func (d *Device) EndianLittle() bool {
	val, _ := d.getInfoBool(C.CL_DEVICE_ENDIAN_LITTLE, true)
	return val
}

// Is CL_TRUE if the device implements error correction for all
// accesses to compute device memory (global and constant). Is
// CL_FALSE if the device does not implement such error correction.
func (d *Device) ErrorCorrectionSupport() bool {
	val, _ := d.getInfoBool(C.CL_DEVICE_ERROR_CORRECTION_SUPPORT, true)
	return val
}

func (d *Device) ImageSupport() bool {
	val, _ := d.getInfoBool(C.CL_DEVICE_IMAGE_SUPPORT, true)
	return val
}

// Is CL_TRUE if the device's preference is for the user to be
// responsible for synchronization, when sharing memory objects
// between OpenCL and other APIs such as DirectX, CL_FALSE if the
// device / implementation has a performant path for performing
// synchronization of memory object shared between OpenCL and other
// APIs such as DirectX
func (d *Device) PreferredInteropUserSync() bool {
	val, _ := d.getInfoBool(C.CL_DEVICE_PREFERRED_INTEROP_USER_SYNC, true)
	return val
}

func (d *Device) Type() DeviceType {
	var deviceType C.cl_device_type
	if err := C.clGetDeviceInfo(d.id, C.CL_DEVICE_TYPE, C.size_t(unsafe.Sizeof(deviceType)), unsafe.Pointer(&deviceType), nil); err != C.CL_SUCCESS {
		panic("Failed to get device type")
	}
	return DeviceType(deviceType)
}

// Describes double precision floating-point capability of the OpenCL device
func (d *Device) DoubleFPConfig() FPConfig {
	var fpConfig C.cl_device_fp_config
	if err := C.clGetDeviceInfo(d.nullableId(), C.CL_DEVICE_DOUBLE_FP_CONFIG, C.size_t(unsafe.Sizeof(fpConfig)), unsafe.Pointer(&fpConfig), nil); err != C.CL_SUCCESS {
		panic("Failed to get double FP config")
	}
	return FPConfig(fpConfig)
}

// Describes single precision floating-point capability of the OpenCL device
func (d *Device) SingleFPConfig() FPConfig {
	var fpConfig C.cl_device_fp_config
	if err := C.clGetDeviceInfo(d.nullableId(), C.CL_DEVICE_SINGLE_FP_CONFIG, C.size_t(unsafe.Sizeof(fpConfig)), unsafe.Pointer(&fpConfig), nil); err != C.CL_SUCCESS {
		panic("Failed to get single FP config")
	}
	return FPConfig(fpConfig)
}

// Describes the OPTIONAL half precision floating-point capability of the OpenCL device
func (d *Device) HalfFPConfig() FPConfig {
	var fpConfig C.cl_device_fp_config
	err := C.clGetDeviceInfo(d.id, C.CL_DEVICE_HALF_FP_CONFIG, C.size_t(unsafe.Sizeof(fpConfig)), unsafe.Pointer(&fpConfig), nil)
	if err != C.CL_SUCCESS {
		return FPConfig(0)
	}
	return FPConfig(fpConfig)
}

// Type of local memory supported. This can be set to CL_LOCAL implying dedicated
// local memory storage such as SRAM, or CL_GLOBAL. For custom devices, CL_NONE
// can also be returned indicating no local memory support.
func (d *Device) LocalMemType() LocalMemType {
	var memType C.cl_device_local_mem_type
	if err := C.clGetDeviceInfo(d.id, C.CL_DEVICE_LOCAL_MEM_TYPE, C.size_t(unsafe.Sizeof(memType)), unsafe.Pointer(&memType), nil); err != C.CL_SUCCESS {
		return LocalMemType(C.CL_NONE)
	}
	return LocalMemType(memType)
}

// Describes the execution capabilities of the device. The mandated minimum capability is CL_EXEC_KERNEL.
func (d *Device) ExecutionCapabilities() ExecCapability {
	var execCap C.cl_device_exec_capabilities
	if err := C.clGetDeviceInfo(d.id, C.CL_DEVICE_EXECUTION_CAPABILITIES, C.size_t(unsafe.Sizeof(execCap)), unsafe.Pointer(&execCap), nil); err != C.CL_SUCCESS {
		panic("Failed to get execution capabilities")
	}
	return ExecCapability(execCap)
}

func (d *Device) GlobalMemCacheType() MemCacheType {
	var memType C.cl_device_mem_cache_type
	if err := C.clGetDeviceInfo(d.nullableId(), C.CL_DEVICE_GLOBAL_MEM_CACHE_TYPE, C.size_t(unsafe.Sizeof(memType)), unsafe.Pointer(&memType), nil); err != C.CL_SUCCESS {
		return MemCacheType(C.CL_NONE)
	}
	return MemCacheType(memType)
}

// Maximum number of work-items that can be specified in each dimension of the work-group to clEnqueueNDRangeKernel.
//
// Returns n size_t entries, where n is the value returned by the query for CL_DEVICE_MAX_WORK_ITEM_DIMENSIONS.
//
// The minimum value is (1, 1, 1) for devices that are not of type CL_DEVICE_TYPE_CUSTOM.
func (d *Device) MaxWorkItemSizes() []int {
	dims := d.MaxWorkItemDimensions()
	sizes := make([]C.size_t, dims)
	if err := C.clGetDeviceInfo(d.nullableId(), C.CL_DEVICE_MAX_WORK_ITEM_SIZES, C.size_t(int(unsafe.Sizeof(sizes[0]))*dims), unsafe.Pointer(&sizes[0]), nil); err != C.CL_SUCCESS {
		panic("Failed to get max work item sizes")
	}
	intSizes := make([]int, dims)
	for i, s := range sizes {
		intSizes[i] = int(s)
	}
	return intSizes
}

// Native vector width size for built-in char type that can be put into vectors.
// The vector width is defined as the number of scalar elements that can be stored in
// the vector.
func (d *Device) NativeVectorWidthChar() int {
	val, _ := d.getInfoUint(C.CL_DEVICE_NATIVE_VECTOR_WIDTH_CHAR, true)
	return int(val)
}

// Native vector width size for built-in short type that can be put into vectors.
// The vector width is defined as the number of scalar elements that can be stored in
// the vector.
func (d *Device) NativeVectorWidthShort() int {
	val, _ := d.getInfoUint(C.CL_DEVICE_NATIVE_VECTOR_WIDTH_SHORT, true)
	return int(val)
}

// Native vector width size for built-in int type that can be put into vectors.
// The vector width is defined as the number of scalar elements that can be stored in
// the vector.
func (d *Device) NativeVectorWidthInt() int {
	val, _ := d.getInfoUint(C.CL_DEVICE_NATIVE_VECTOR_WIDTH_INT, true)
	return int(val)
}

// Native vector width size for built-in long type that can be put into vectors.
// The vector width is defined as the number of scalar elements that can be stored in
// the vector.
func (d *Device) NativeVectorWidthLong() int {
	val, _ := d.getInfoUint(C.CL_DEVICE_NATIVE_VECTOR_WIDTH_LONG, true)
	return int(val)
}

// Native vector width size for built-in float type that can be put into vectors.
// The vector width is defined as the number of scalar elements that can be stored in
// the vector.
func (d *Device) NativeVectorWidthFloat() int {
	val, _ := d.getInfoUint(C.CL_DEVICE_NATIVE_VECTOR_WIDTH_FLOAT, true)
	return int(val)
}

// Native vector width size for built-in double type that can be put into vectors.
// The vector width is defined as the number of scalar elements that can be stored in
// the vector. Must return 0 when cl_khr_fp64 is unsupported.
func (d *Device) NativeVectorWidthDouble() int {
	val, _ := d.getInfoUint(C.CL_DEVICE_NATIVE_VECTOR_WIDTH_DOUBLE, true)
	return int(val)
}

// Native vector width size for built-in half type that can be put into vectors.
// The vector width is defined as the number of scalar elements that can be stored in
// the vector. Must return 0 when cl_khr_fp16 is unsupported.
func (d *Device) NativeVectorWidthHalf() int {
	val, _ := d.getInfoUint(C.CL_DEVICE_NATIVE_VECTOR_WIDTH_HALF, true)
	return int(val)
}

// Preferred native vector width size for built-in char type that can be put into vectors.
// The vector width is defined as the number of scalar elements that can be stored in
// the vector.
func (d *Device) PreferredVectorWidthChar() int {
	val, _ := d.getInfoUint(C.CL_DEVICE_PREFERRED_VECTOR_WIDTH_CHAR, true)
	return int(val)
}

// Preferred native vector width size for built-in short type that can be put into vectors.
// The vector width is defined as the number of scalar elements that can be stored in
// the vector.
func (d *Device) PreferredVectorWidthShort() int {
	val, _ := d.getInfoUint(C.CL_DEVICE_PREFERRED_VECTOR_WIDTH_SHORT, true)
	return int(val)
}

// Preferred native vector width size for built-in int type that can be put into vectors.
// The vector width is defined as the number of scalar elements that can be stored in
// the vector.
func (d *Device) PreferredVectorWidthInt() int {
	val, _ := d.getInfoUint(C.CL_DEVICE_PREFERRED_VECTOR_WIDTH_INT, true)
	return int(val)
}

// Preferred native vector width size for built-in long type that can be put into vectors.
// The vector width is defined as the number of scalar elements that can be stored in
// the vector.
func (d *Device) PreferredVectorWidthLong() int {
	val, _ := d.getInfoUint(C.CL_DEVICE_PREFERRED_VECTOR_WIDTH_LONG, true)
	return int(val)
}

// Preferred native vector width size for built-in float type that can be put into vectors.
// The vector width is defined as the number of scalar elements that can be stored in
// the vector.
func (d *Device) PreferredVectorWidthFloat() int {
	val, _ := d.getInfoUint(C.CL_DEVICE_PREFERRED_VECTOR_WIDTH_FLOAT, true)
	return int(val)
}

// Preferred native vector width size for built-in double type that can be put into vectors.
// The vector width is defined as the number of scalar elements that can be stored in
// the vector. Must return 0 when cl_khr_fp64 is unsupported.
func (d *Device) PreferredVectorWidthDouble() int {
	val, _ := d.getInfoUint(C.CL_DEVICE_PREFERRED_VECTOR_WIDTH_DOUBLE, true)
	return int(val)
}

// Preferred native vector width size for built-in double type that can be put into vectors.
// The vector width is defined as the number of scalar elements that can be stored in
// the vector. Must return 0 when cl_khr_fp16 is unsupported.
func (d *Device) PreferredVectorWidthHalf() int {
	val, _ := d.getInfoUint(C.CL_DEVICE_PREFERRED_VECTOR_WIDTH_HALF, true)
	return int(val)
}

func (d *Device) QueueProperties() CommandQueueProperty {
	var val C.cl_command_queue_properties
	if err := C.clGetDeviceInfo(d.nullableId(), C.CL_DEVICE_QUEUE_PROPERTIES, C.size_t(unsafe.Sizeof(val)), unsafe.Pointer(&val), nil); err != C.CL_SUCCESS {
		panic("Should never fail")
		return 0
	}
	return CommandQueueProperty(val)
}

func (d *Device) PartitionDeviceEqually(n int) ([]*Device, error) {
	var deviceList []C.cl_device_id
	var deviceCount C.cl_uint
	defer C.free(unsafe.Pointer(&deviceList))
	defer C.free(unsafe.Pointer(&deviceCount))
	err := C.clCreateSubDevices(d.nullableId(), C.partitionDeviceEqually((C.uint)(n)), 1, &deviceList[0], &deviceCount)
	if toError(err) != nil {
		return nil, toError(err)
	}
	val := make([]*Device, int(deviceCount))
	for idx := range val {
		val[idx].id = deviceList[idx]
	}
	return val, nil
}

func (d *Device) PartitionDeviceByCounts(n []int) ([]*Device, error) {
	var deviceList []C.cl_device_id
	var deviceCount C.cl_uint
	defer C.free(unsafe.Pointer(&deviceList))
	defer C.free(unsafe.Pointer(&deviceCount))

	Counts := make([]C.uint, len(n))
	defer C.free(unsafe.Pointer(&Counts))
	for ii, nn := range n {
		Counts[ii] = (C.uint)(nn)
	}
	err := C.clCreateSubDevices(d.nullableId(), C.partitionDeviceByCounts(&Counts[0], (C.uint)(len(n))), 1, &deviceList[0], &deviceCount)
	if toError(err) != nil {
		return nil, toError(err)
	}
	val := make([]*Device, int(deviceCount))
	for idx := range val {
		val[idx].id = deviceList[idx]
	}
	return val, nil
}

func (d *Device) PartitionDeviceByNumaDomain(n []int) ([]*Device, error) {
	var deviceList []C.cl_device_id
	var deviceCount C.cl_uint
	defer C.free(unsafe.Pointer(&deviceList))
	defer C.free(unsafe.Pointer(&deviceCount))

	Counts := make([]C.uint, len(n))
	defer C.free(unsafe.Pointer(&Counts))
	for ii, nn := range n {
		Counts[ii] = (C.uint)(nn)
	}
	err := C.clCreateSubDevices(d.nullableId(), C.partitionDeviceByNuma(), 1, &deviceList[0], &deviceCount)
	if toError(err) != nil {
		return nil, toError(err)
	}
	val := make([]*Device, int(deviceCount))
	for idx := range val {
		val[idx].id = deviceList[idx]
	}
	return val, nil
}

func (d *Device) PartitionDeviceByL4CacheDomain(n []int) ([]*Device, error) {
	var deviceList []C.cl_device_id
	var deviceCount C.cl_uint
	defer C.free(unsafe.Pointer(&deviceList))
	defer C.free(unsafe.Pointer(&deviceCount))

	Counts := make([]C.uint, len(n))
	defer C.free(unsafe.Pointer(&Counts))
	for ii, nn := range n {
		Counts[ii] = (C.uint)(nn)
	}
	err := C.clCreateSubDevices(d.nullableId(), C.partitionDeviceByL4Cache(), 1, &deviceList[0], &deviceCount)
	if toError(err) != nil {
		return nil, toError(err)
	}
	val := make([]*Device, int(deviceCount))
	for idx := range val {
		val[idx].id = deviceList[idx]
	}
	return val, nil
}

func (d *Device) PartitionDeviceByL3CacheDomain(n []int) ([]*Device, error) {
	var deviceList []C.cl_device_id
	var deviceCount C.cl_uint
	defer C.free(unsafe.Pointer(&deviceList))
	defer C.free(unsafe.Pointer(&deviceCount))

	Counts := make([]C.uint, len(n))
	defer C.free(unsafe.Pointer(&Counts))
	for ii, nn := range n {
		Counts[ii] = (C.uint)(nn)
	}
	err := C.clCreateSubDevices(d.nullableId(), C.partitionDeviceByL3Cache(), 1, &deviceList[0], &deviceCount)
	if toError(err) != nil {
		return nil, toError(err)
	}
	val := make([]*Device, int(deviceCount))
	for idx := range val {
		val[idx].id = deviceList[idx]
	}
	return val, nil
}

func (d *Device) PartitionDeviceByL2CacheDomain(n []int) ([]*Device, error) {
	var deviceList []C.cl_device_id
	var deviceCount C.cl_uint
	defer C.free(unsafe.Pointer(&deviceList))
	defer C.free(unsafe.Pointer(&deviceCount))

	Counts := make([]C.uint, len(n))
	defer C.free(unsafe.Pointer(&Counts))
	for ii, nn := range n {
		Counts[ii] = (C.uint)(nn)
	}
	err := C.clCreateSubDevices(d.nullableId(), C.partitionDeviceByL2Cache(), 1, &deviceList[0], &deviceCount)
	if toError(err) != nil {
		return nil, toError(err)
	}
	val := make([]*Device, int(deviceCount))
	for idx := range val {
		val[idx].id = deviceList[idx]
	}
	return val, nil
}

func (d *Device) PartitionDeviceByL1CacheDomain(n []int) ([]*Device, error) {
	var deviceList []C.cl_device_id
	var deviceCount C.cl_uint
	defer C.free(unsafe.Pointer(&deviceList))
	defer C.free(unsafe.Pointer(&deviceCount))

	Counts := make([]C.uint, len(n))
	defer C.free(unsafe.Pointer(&Counts))
	for ii, nn := range n {
		Counts[ii] = (C.uint)(nn)
	}
	err := C.clCreateSubDevices(d.nullableId(), C.partitionDeviceByL1Cache(), 1, &deviceList[0], &deviceCount)
	if toError(err) != nil {
		return nil, toError(err)
	}
	val := make([]*Device, int(deviceCount))
	for idx := range val {
		val[idx].id = deviceList[idx]
	}
	return val, nil
}

func (d *Device) PartitionDeviceByNextPartitionableDomain(n []int) ([]*Device, error) {
	var deviceList []C.cl_device_id
	var deviceCount C.cl_uint
	defer C.free(unsafe.Pointer(&deviceList))
	defer C.free(unsafe.Pointer(&deviceCount))

	Counts := make([]C.uint, len(n))
	defer C.free(unsafe.Pointer(&Counts))
	for ii, nn := range n {
		Counts[ii] = (C.uint)(nn)
	}
	err := C.clCreateSubDevices(d.nullableId(), C.partitionDeviceByNextPartitionable(), 1, &deviceList[0], &deviceCount)
	if toError(err) != nil {
		return nil, toError(err)
	}
	val := make([]*Device, int(deviceCount))
	for idx := range val {
		val[idx].id = deviceList[idx]
	}
	return val, nil
}

func (d *Device) PartitionAffinityDomain() DeviceAffinityDomain {
	var deviceAffinityDomain C.cl_device_affinity_domain
	var paramSize C.size_t
	defer C.free(unsafe.Pointer(&deviceAffinityDomain))
	defer C.free(unsafe.Pointer(&paramSize))
	if err := C.CLGetDeviceInfoParamSize(d.nullableId(), C.CL_DEVICE_PARTITION_AFFINITY_DOMAIN, &paramSize); err != C.CL_SUCCESS {
		panic("Should never fail getting parameter size for device info")
	}
	if err := C.CLGetDeviceInfoParamUnsafe(d.nullableId(), C.CL_DEVICE_PARTITION_AFFINITY_DOMAIN, paramSize, unsafe.Pointer(&deviceAffinityDomain)); err != C.CL_SUCCESS {
		panic("Should never fail getting device info")
	}
	res := DeviceAffinityDomain(deviceAffinityDomain)
	return res
}

func (d *Device) ParentDevice() *Device {
	var devId C.cl_device_id
	var paramSize C.size_t
	defer C.free(unsafe.Pointer(&devId))
	defer C.free(unsafe.Pointer(&paramSize))
	if err := C.CLGetDeviceInfoParamSize(d.nullableId(), C.CL_DEVICE_PARENT_DEVICE, &paramSize); err != C.CL_SUCCESS {
		panic("Should never fail getting parameter size for device info")
	}
	if err := C.CLGetDeviceInfoParamUnsafe(d.nullableId(), C.CL_DEVICE_PARENT_DEVICE, paramSize, unsafe.Pointer(&devId)); err != C.CL_SUCCESS {
		panic("Should never fail getting device info")
	}
	res := new(Device)
	res.id = devId
	return res
}
