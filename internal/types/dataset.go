package types

type Dataset string

const (
	DatasetMHYClusters Dataset = "mgy_clusters_2022_05.fa"
	DatasetBFD         Dataset = "bfd-first_non_consensus_sequences.fasta"
	DatasetUniRef90    Dataset = "uniref90_2022_05.fa"
	DatasetPDB         Dataset = "pdb_2022_09_28_mmcif_files.tar"
	DatasetPDBSeqReq   Dataset = "pdb_seqres_2022_09_28.fasta"
	DatasetRNACentral  Dataset = "rnacentral_active_seq_id_90_cov_80_linclust.fasta"
	DatasetNT          Dataset = "nt_rna_2023_02_23_clust_seq_id_90_cov_80_rep_seq.fasta"
	DatasetRFam        Dataset = "rfam_14_9_clust_seq_id_90_cov_80_rep_seq.fasta"
)

var Datasets = []Dataset{
	DatasetMHYClusters,
	DatasetBFD,
	DatasetUniRef90,
	DatasetPDB,
	DatasetPDBSeqReq,
	DatasetRNACentral,
	DatasetNT,
	DatasetRFam,
}

func (d Dataset) String() string {
	return string(d)
}

func (d Dataset) Size() int64 {
	switch d {
	case DatasetRFam:
		return 228433680
	default:
		return 0
	}
}

func (d Dataset) Hash() string {
	switch d {
	case DatasetRFam:
		return "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
	default:
		return ""
	}
}
