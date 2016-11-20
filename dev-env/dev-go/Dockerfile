FROM centos:centos7.0.1406

# Install git, gcc, vim, ctags, cscope
RUN yum update -y && \
    yum install -y git && \
    yum install -y gcc && \
    yum install -y wget && \
    yum install -y ctags && \
    yum install -y cscope && \
    yum install -y tar

# Install golang version 1.7.3
RUN cd /home && \
    wget -q https://storage.googleapis.com/golang/go1.7.3.linux-amd64.tar.gz && \
    tar -C /usr/local -xzf /home/go1.7.3.linux-amd64.tar.gz && \
    rm *.gz && \
    mkdir -p /go

# Configure go env variables
ENV GOPATH=/go
ENV PATH=$PATH:/usr/local/go/bin
ENV PATH=$PATH:$GOPATH/bin

# Copy local vimrc to root
COPY .vimrc /root/
# Install vim and then configure vimrc and vim plugins
RUN yum install -y vim && \

    # Install Pathogen
    yum install -y curl && \
    mkdir -p ~/.vim/autoload ~/.vim/bundle && \
    curl -LSso ~/.vim/autoload/pathogen.vim https://tpo.pe/pathogen.vim && \
    echo "execute pathogen#infect()" >> ~/.vimrc && \
    echo "filetype plugin indent on" >> ~/.vimrc && \

    # Install vim-go
    git clone https://github.com/fatih/vim-go.git ~/.vim/bundle/vim-go && \
    vim +GoInstallBinaries +qall && \
    vim +Helptags +qall && \

    # Install YouCompleteMe for auto completion of golang
    git clone https://github.com/Valloric/YouCompleteMe.git ~/.vim/bundle/YouCompleteMe && \
    yum install -y build-essential cmake && \
    yum install -y python-devel && \
    yum install -y python && \
   # cd ~/.vim/bundle/YouCompleteMe && \
   # git submodule update --init --recursive && \
   # ~/.vim/bundle/YouCompleteMe/install.py --gocode-completer && \

    # Install tagbar to show function info
    git clone https://github.com/majutsushi/tagbar.git ~/.vim/bundle/tagbar && \
    echo "nmap <F8> :TagbarToggle<CR>" >> ~/.vimrc && \

    # Install g-explorer for go docs
    git clone https://github.com/garyburd/go-explorer.git ~/.vim/bundle/go-explorer && \

    # Enable syntax highlighting for go
    (echo "let g:go_highlight_functions = 1"; \
     echo "let g:go_highlight_methods = 1"; \
     echo "let g:go_highlight_structs = 1"; \
     echo "let g:go_highlight_operators = 1"; \
     echo "let g:go_highlight_build_constraints = 1") \
    >> ~/.vimrc && \

    # Install NERDTree plugin
    git clone https://github.com/scrooloose/nerdtree.git ~/.vim/bundle/nerdtree && \
    echo "map <C-n> :NERDTreeToggle<CR>" >> ~/.vimrc && \
    echo "set encoding=utf-8" && \

    # Install vim-airline
    git clone https://github.com/vim-airline/vim-airline ~/.vim/bundle/vim-airline && \
    git clone https://github.com/vim-airline/vim-airline-themes ~/.vim/bundle/vim-airline-themes && \

    # Add cscope mappings
    mkdir ~/.vim/plugin && \
    cd ~/.vim/plugin && \
    wget -q http://cscope.sourceforge.net/cscope_maps.vim

WORKDIR /root
