# Coddy Sandbox Image
# Provides isolated execution environment for Python and Node.js code

FROM python:3.12-slim-bookworm

# Set environment variables
ENV PYTHONDONTWRITEBYTECODE=1
ENV PYTHONUNBUFFERED=1
ENV NODE_VERSION=20.x
ENV USER=coddy
ENV HOME=/home/$USER

# Install system dependencies
RUN apt-get update && apt-get install -y --no-install-recommends \
    curl \
    ca-certificates \
    gcc \
    g++ \
    make \
    git \
    # Required for matplotlib
    libfreetype6-dev \
    libpng-dev \
    # Required for Pillow
    libjpeg-dev \
    libtiff-dev \
    # Required for scipy/numpy
    libopenblas-dev \
    gfortran \
    # Required for Node.js build tools
    build-essential \
    && rm -rf /var/lib/apt/lists/*

# Install Node.js
RUN curl -fsSL https://deb.nodesource.com/setup_${NODE_VERSION} | bash - \
    && apt-get install -y nodejs \
    && rm -rf /var/lib/apt/lists/*

# Create non-root user
RUN useradd -m -s /bin/bash $USER

# Switch to user
USER $USER
WORKDIR $HOME

# Install Python packages
RUN pip install --user --no-cache-dir \
    numpy \
    pandas \
    matplotlib \
    scipy \
    scikit-learn \
    requests \
    beautifulsoup4 \
    Pillow \
    sympy \
    openpyxl \
    ipython

# Install Node.js packages globally (in user directory)
RUN npm config set prefix $HOME/.local && \
    npm install -g \
    lodash \
    axios \
    cheerio \
    csv-parse \
    mathjs \
    sharp

# Ensure local bin is in PATH
ENV PATH=$HOME/.local/bin:$PATH

# Create working directory for code execution
RUN mkdir -p $HOME/workspace
WORKDIR $HOME/workspace

# Metadata
LABEL maintainer="Coddy Contributors"
LABEL description="Coddy Sandbox - Isolated execution environment for LLM code execution"
LABEL version="0.1.0"

# Default command (container should be kept alive by orchestrator)
CMD ["tail", "-f", "/dev/null"]
