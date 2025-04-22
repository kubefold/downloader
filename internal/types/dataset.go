package types

type Dataset string

const (
	DatasetMHYClusters Dataset = "mgy_clusters_2022_05.fa"
	DatasetBFD         Dataset = "bfd-first_non_consensus_sequences.fasta"
	DatasetUniRef90    Dataset = "uniref90_2022_05.fa"
	DatasetUniProt     Dataset = "uniprot_all_2021_04.fa"
	DatasetPDB         Dataset = "pdb_2022_09_28_mmcif_files.tar"
	DatasetPDBSeqReq   Dataset = "pdb_seqres_2022_09_28.fasta"
	DatasetRNACentral  Dataset = "rnacentral_active_seq_id_90_cov_80_linclust.fasta"
	DatasetNT          Dataset = "nt_rna_2023_02_23_clust_seq_id_90_cov_80_rep_seq.fasta"
	DatasetRFam        Dataset = "rfam_14_9_clust_seq_id_90_cov_80_rep_seq.fasta"
)

var Datasets = []Dataset{
	DatasetMHYClusters,
	DatasetBFD,
	DatasetUniProt,
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
	case DatasetBFD:
		return 18171626364
	case DatasetMHYClusters:
		return 128579703018
	case DatasetNT:
		return 80977012680
	case DatasetPDB:
		return 250521374720
	case DatasetPDBSeqReq:
		return 232899463
	case DatasetUniProt:
		return 108447942931
	case DatasetUniRef90:
		return 71821260491
	case DatasetRNACentral:
		return 13860314914
	default:
		return 0
	}
}

func (d Dataset) Hash() string {
	switch d {
	case DatasetRFam:
		return "55ef718071244ad7433678ba249aaeb67707b499f0189a38edadca8d64972318"
	case DatasetBFD:
		return "fd87dca06401b03f4ac3c59a82dac14db491a7933ed6abaa19e14e02c6eb1af5"
	case DatasetMHYClusters:
		return "9e7f50956c19cbcd8181dc5e9d7d6eebc08257cc858fc07d3ec88fd6b48dbbc9"
	case DatasetNT:
		return "14c05ac0827c9bf06a37acfc4b3dd1d66e48d5a5f713c0de68611aa7fedc00f9"
	case DatasetPDBSeqReq:
		return "1b3bc853322c32f2eea818065b8f569a18d25a52326a8d2c2c3de85752e55fe1"
	case DatasetUniProt:
		return "76f32efd5c6ba73857b0beb3bf1ff823cf0dbef3d876c70d80ee387db13a169d"
	case DatasetUniRef90:
		return "f0c61e13a6f71ec2b19e44d35acb531ed3a06a4a839fc12feb80d3adf883c049"
	case DatasetRNACentral:
		return "6c33f15c48d2ac8d7d42a8699ff2e7bd6a4816f8a074157522d3c5b591f927eb"
	default:
		return ""
	}
}
