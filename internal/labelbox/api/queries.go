package api

var projectDatasetRowsQuery = `
	query datasetRows(
		$whereProject: WhereUniqueIdInput!,
		$whereDatasets: DatasetWhereInput,
		$skip: Int,
		$pageSizeDatarows: PageSize
	) {
		project(where: $whereProject) {
			datasets(where: $whereDatasets) {
				id
				dataRows(skip: $skip, first:$pageSizeDatarows) {
					id
					externalId
				}
			}
		}
	}
`
