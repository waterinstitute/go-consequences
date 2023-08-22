FROM golang:1.20.7

# Install GDAL dependencies
RUN apt-get update && \
    apt-get install -y --no-install-recommends \
    libgdal-dev \
    gdal-bin \
    wget \
    cmake \
    libcurl4-gnutls-dev

# Set GDAL version
#ENV GDAL_VERSION=3.7.0
ENV LD_LIBRARY_PATH=/usr/local/lib
ENV SI_VERSION="0_1"
# Download and build GDAL from source
#RUN wget https://download.osgeo.org/gdal/${GDAL_VERSION}/gdal-${GDAL_VERSION}.tar.gz && \
#    tar -xzf gdal-${GDAL_VERSION}.tar.gz && \
#    cd gdal-${GDAL_VERSION} && \
#    mkdir build && \
#    cd build && \
#    cmake .. && \
#    cmake --build . && \
#    cmake --build . --target install
# Set the working directory
RUN mkdir app
WORKDIR /app
COPY main /app/
COPY structure_inventory/ /app/structure_inventory/
# Set the entrypoint command
ENTRYPOINT ["./main","-raster=raster.tif","-result=result.shp","-si=si.csv"]
