FROM golang:1.21-bullseye

RUN chmod -R 777 /usr/local/go/bin

ARG USERNAME=vscode
RUN useradd -m -s /bin/zsh $USERNAME
USER $USERNAME

RUN go install go.opentelemetry.io/collector/cmd/builder@v0.91.0
RUN go install github.com/open-telemetry/opentelemetry-collector-contrib/cmd/mdatagen@v0.91.0
RUN go install golang.org/x/tools/gopls@latest
RUN go install github.com/go-delve/delve/cmd/dlv@latest

# set up oh-my-zsh
RUN sh -c "$(wget -O- https://raw.githubusercontent.com/robbyrussell/oh-my-zsh/master/tools/install.sh)" "" --unattended &> /dev/null
ENV ENV="/home/$USERNAME/.ashrc" \
    ZSH=/home/$USERNAME/.oh-my-zsh \
    EDITOR=vi \
    LANG=en_US.UTF-8
RUN printf 'ZSH_THEME="robbyrussell"\nENABLE_CORRECTION="false"\nplugins=(git copyfile extract colorize dotenv encode64 golang)\nsource $ZSH/oh-my-zsh.sh\nexport PATH=$PATH:~/.local/bin' > "/home/$USERNAME/.zshrc"
RUN echo "exec `which zsh`" > "/home/$USERNAME/.ashrc"