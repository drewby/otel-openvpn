FROM golang:1.21-bullseye

ARG USERNAME=vscode
RUN useradd -m -s /bin/zsh $USERNAME
USER $USERNAME

RUN go install go.opentelemetry.io/collector/cmd/builder@latest
RUN go install github.com/open-telemetry/opentelemetry-collector-contrib/cmd/mdatagen@latest
RUN go install golang.org/x/tools/gopls@latest

# set up oh-my-zsh
RUN sh -c "$(wget -O- https://raw.githubusercontent.com/robbyrussell/oh-my-zsh/master/tools/install.sh)" "" --unattended &> /dev/null
ENV ENV="/home/$USERNAME/.ashrc" \
    ZSH=/home/$USERNAME/.oh-my-zsh \
    EDITOR=vi \
    LANG=en_US.UTF-8
RUN printf 'ZSH_THEME="robbyrussell"\nENABLE_CORRECTION="false"\nplugins=(git copyfile extract colorize dotenv encode64 golang)\nsource $ZSH/oh-my-zsh.sh\nexport PATH=$PATH:~/.local/bin' > "/home/$USERNAME/.zshrc"
RUN echo "exec `which zsh`" > "/home/$USERNAME/.ashrc"