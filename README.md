# 🧬 Protein Database Downloader

A tool for downloading and decompressing protein databases for bioinformatics research.

## About

Protein Database Downloader is a utility designed to fetch large protein datasets from AlphaFold's database. It supports rate limiting, progress tracking, and automatic decompression of downloaded files.

Developed by Mateusz Woźniak <matisiek11@gmail.com>

## Supported Datasets

The following protein datasets are supported:

- `mgy_clusters_2022_05.fa` - MGY Clusters
- `bfd-first_non_consensus_sequences.fasta` - BFD non-consensus sequences
- `uniref90_2022_05.fa` - UniRef90
- `uniprot_all_2021_04.fa` - UniProt
- `pdb_2022_09_28_mmcif_files.tar` - PDB mmCIF files
- `pdb_seqres_2022_09_28.fasta` - PDB sequence resources
- `rnacentral_active_seq_id_90_cov_80_linclust.fasta` - RNACentral
- `nt_rna_2023_02_23_clust_seq_id_90_cov_80_rep_seq.fasta` - NT RNA
- `rfam_14_9_clust_seq_id_90_cov_80_rep_seq.fasta` - RFam

## Usage

### Environment Variables

The application is configured using the following environment variables:

- `DATASET` (required): The dataset to download (must be one of the supported datasets listed above)
- `DESTINATION` (required): The directory path where the downloaded dataset will be saved
- `RATE` (optional): Download rate limit in KB/s (default: unlimited)

### Docker

```bash
docker run -e DATASET=rfam_14_9_clust_seq_id_90_cov_80_rep_seq.fasta \
           -e DESTINATION=/data \
           -e RATE=1024 \
           -v /local/path:/data \
           kubefold/downloader
```

### Building from Source

```bash
git clone https://github.com/kubefold/downloader.git
cd downloader
go build -o downloader ./cmd/main.go
```

### Running the Binary

```bash
DATASET=rfam_14_9_clust_seq_id_90_cov_80_rep_seq.fasta \
DESTINATION=/data \
RATE=1024 \
./downloader
```

## Features

- Downloads datasets from AlphaFold's database
- Automatically decompresses zstd-compressed files
- Configurable download rate limiting
- Progress tracking with size information
- SHA-256 hash verification for downloaded files
