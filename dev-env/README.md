# Docker images for cns project
It consists of following images:
* dev-go:
  * Ubuntu 14.04
  * git, gcc, vim, wget, ctags, curl
  * vim plugins: vim-go, YCM (auto-completion), tagbar, go-explorer
    * vim-go: https://github.com/fatih/vim-go
    * YCM: https://github.com/fatih/vim-go.git
	* tagbar: https://github.com/majutsushi/tagbar.git
	* g-explorer: https://github.com/garyburd/go-explorer.git
	* NERDTree: https://github.com/scrooloose/nerdtree.git
	* vim-airline: https://github.com/vim-airline/vim-airline
	* Some other syntax highlighting
	
* dev-go-fuse:
  * dev-go plus go-fuse package: https://github.com/bazil/fuse
* dev-go-ceph:
  * dev-go plus go-ceph: https://github.com/noahdesu/go-ceph
* dev-fuse-ceph:
  * dev-go-fuse plus go-ceph


Always use latest images. Images are also uploaded to https://secure-registry.gsintlab.com.

### How to run Docker container.

Map /cns/core/host/go --> /go/src/apporbit inside container. <br />

Example: "docker run -it -v /cns/core/host/go:/go/src/apporbit secure-registry.gsintlab.com/dev-go bash"
