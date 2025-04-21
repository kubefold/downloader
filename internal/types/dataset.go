package types

type Dataset string

const (
	DatasetMHYClusters Dataset = "mgy_clusters_2022_05.fa.zst"
	DatasetBFD         Dataset = "bfd-first_non_consensus_sequences.fasta.zst"
	DatasetUniRef90    Dataset = "uniref90_2022_05.fa uniprot_all_2021_04.fa.zst"
	DatasetPDB         Dataset = "pdb_2022_09_28_mmcif_files.tar.zst"
	DatasetPDBSeqReq   Dataset = "pdb_seqres_2022_09_28.fasta.zst"
	DatasetRNACentral  Dataset = "rnacentral_active_seq_id_90_cov_80_linclust.fasta.zst"
	DatasetNT          Dataset = "nt_rna_2023_02_23_clust_seq_id_90_cov_80_rep_seq.fasta.zst"
	DatasetRFam        Dataset = "rfam_14_9_clust_seq_id_90_cov_80_rep_seq.fasta.zst"
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
