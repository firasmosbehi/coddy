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
    libfreetype6-dev \
    libpng-dev \
    libjpeg-dev \
    libtiff-dev \
    libopenblas-dev \
    gfortran \
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
    openpyxl

# Install Node.js packages
RUN npm config set prefix $HOME/.local && \
    npm install -g \
    lodash \
    axios \
    cheerio \
    csv-parse \
    mathjs

# Ensure local bin is in PATH
ENV PATH=$HOME/.local/bin:$PATH

# Create working directory
RUN mkdir -p $HOME/workspace
WORKDIR $HOME/workspace

# Default command
CMD ["tail", "-f", "/dev/null"]
