# go2opencl
Go bindings to OpenCL version 1.2

The Go bindings to OpenCL use the headers available through Khronos Group.
The current bindings use OpenCL 1.2, which provides the best compatibility
across GPUs from different vendors. OpenCL 3.0 is essentially version 1.2
as of this writing.

The codes were originally tested with Go 1.9 but we set the requirement to
be Go 1.12.

Bindings for VkFFT (https://github.com/DTolm/VkFFT) are included for
convenience.

- D. Tolmachev, IEEE Access vol. 11, pp. 12039-12058
  doi:10.1109/ACCESS.2023.3242240.
