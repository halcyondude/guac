# GUAC Schema Documentation

## Entities
### artifacts
#### Fields
| Field | Type | Comment |
| --- | --- | --- |
| id | UUID |  |
| algorithm | STRING |  |
| digest | STRING |  |

#### Natural Keys (Composite Unique Constraints)
| Keys |
| --- |
| `[algorithm, digest]` |

### bill_of_materials
#### Fields
| Field | Type | Comment |
| --- | --- | --- |
| id | UUID |  |
| uri | STRING | SBOM's URI |
| algorithm | STRING | Digest algorithm |
| digest | STRING |  |
| download_location | STRING |  |
| origin | STRING |  |
| collector | STRING | GUAC collector for the document |
| document_ref | STRING |  |
| known_since | TIMESTAMP |  |
| included_packages_hash | STRING | An opaque hash of the included packages |
| included_artifacts_hash | STRING | An opaque hash of the included artifacts |
| included_dependencies_hash | STRING | An opaque hash of the included dependencies |
| included_occurrences_hash | STRING | An opaque hash of the included occurrences |
| package_id | UUID |  |
| artifact_id | UUID |  |

#### Edges
| To | Edge Field | On Delete |
| --- | --- | --- |
| package_versions | bill_of_materials_package_versions_package | package_id |
| artifacts | bill_of_materials_artifacts_artifact | artifact_id |
| package_versions | bill_of_materials_included_software_packages |  |
| artifacts | bill_of_materials_included_software_artifacts |  |
| dependencies | bill_of_materials_included_dependencies |  |
| occurrences | bill_of_materials_included_occurrences |  |

#### Natural Keys (Composite Unique Constraints)
| Keys |
| --- |
| `[algorithm, digest, uri, download_location, known_since, included_packages_hash, included_artifacts_hash, included_dependencies_hash, included_occurrences_hash, origin, collector, document_ref, package]` |
| `[algorithm, digest, uri, download_location, known_since, included_packages_hash, included_artifacts_hash, included_dependencies_hash, included_occurrences_hash, origin, collector, document_ref, artifact]` |

### builders
#### Fields
| Field | Type | Comment |
| --- | --- | --- |
| id | UUID |  |
| uri | STRING | The URI of the builder, used as a unique identifier in the graph query |

#### Natural Keys (Composite Unique Constraints)
| Keys |
| --- |
| `[uri]` |

### certifications
#### Fields
| Field | Type | Comment |
| --- | --- | --- |
| id | UUID |  |
| type | STRING |  |
| justification | STRING |  |
| known_since | TIMESTAMP |  |
| origin | STRING |  |
| collector | STRING |  |
| document_ref | STRING |  |
| source_id | UUID |  |
| package_version_id | UUID |  |
| package_name_id | UUID |  |
| artifact_id | UUID |  |

#### Edges
| To | Edge Field | On Delete |
| --- | --- | --- |
| source_names | certifications_source_names_source | source_id |
| package_versions | certifications_package_versions_package_version | package_version_id |
| package_names | certifications_package_names_all_versions | package_name_id |
| artifacts | certifications_artifacts_artifact | artifact_id |

#### Natural Keys (Composite Unique Constraints)
| Keys |
| --- |
| `[type, justification, origin, collector, source, known_since, document_ref]` |
| `[type, justification, origin, collector, package_version, known_since, document_ref]` |
| `[type, justification, origin, collector, package_name, known_since, document_ref]` |
| `[type, justification, origin, collector, artifact, known_since, document_ref]` |

### certify_legals
#### Fields
| Field | Type | Comment |
| --- | --- | --- |
| id | UUID |  |
| declared_license | STRING |  |
| discovered_license | STRING |  |
| attribution | STRING |  |
| justification | STRING |  |
| time_scanned | TIMESTAMP |  |
| origin | STRING |  |
| collector | STRING |  |
| document_ref | STRING |  |
| declared_licenses_hash | STRING | An opaque hash of the declared license IDs to ensure uniqueness |
| discovered_licenses_hash | STRING | An opaque hash of the discovered license IDs to ensure uniqueness |
| package_id | UUID |  |
| source_id | UUID |  |

#### Edges
| To | Edge Field | On Delete |
| --- | --- | --- |
| package_versions | certify_legals_package_versions_package | package_id |
| source_names | certify_legals_source_names_source | source_id |
| licenses | certify_legal_declared_licenses |  |
| licenses | certify_legal_discovered_licenses |  |

#### Natural Keys (Composite Unique Constraints)
| Keys |
| --- |
| `[declared_license, justification, origin, collector, declared_licenses_hash, discovered_licenses_hash, source]` |
| `[declared_license, justification, origin, collector, declared_licenses_hash, discovered_licenses_hash, package]` |

### certify_scorecards
#### Fields
| Field | Type | Comment |
| --- | --- | --- |
| id | UUID |  |
| checks | STRING |  |
| aggregate_score | DOUBLE | Overall Scorecard score for the source |
| time_scanned | TIMESTAMP |  |
| scorecard_version | STRING |  |
| scorecard_commit | STRING |  |
| origin | STRING |  |
| collector | STRING |  |
| document_ref | STRING |  |
| checks_hash | STRING | A SHA1 of the checks fields after sorting keys, used to ensure uniqueness of scorecard records. |
| source_id | UUID |  |

#### Edges
| To | Edge Field | On Delete |
| --- | --- | --- |
| source_names | certify_scorecards_source_names_source | source_id |

#### Natural Keys (Composite Unique Constraints)
| Keys |
| --- |
| `[source, origin, collector, scorecard_version, scorecard_commit, aggregate_score, time_scanned, checks_hash, document_ref]` |

### certify_vexes
#### Fields
| Field | Type | Comment |
| --- | --- | --- |
| id | UUID |  |
| known_since | TIMESTAMP |  |
| status | STRING |  |
| statement | STRING |  |
| status_notes | STRING |  |
| justification | STRING |  |
| origin | STRING |  |
| collector | STRING |  |
| document_ref | STRING |  |
| package_id | UUID |  |
| artifact_id | UUID |  |
| vulnerability_id | UUID |  |

#### Edges
| To | Edge Field | On Delete |
| --- | --- | --- |
| package_versions | certify_vexes_package_versions_package | package_id |
| artifacts | certify_vexes_artifacts_artifact | artifact_id |
| vulnerability_ids | certify_vexes_vulnerability_ids_vulnerability | vulnerability_id |

#### Natural Keys (Composite Unique Constraints)
| Keys |
| --- |
| `[known_since, justification, status, origin, collector, document_ref, vulnerability, package]` |
| `[known_since, justification, status, origin, collector, document_ref, vulnerability, artifact]` |

### certify_vulns
#### Fields
| Field | Type | Comment |
| --- | --- | --- |
| id | UUID |  |
| time_scanned | TIMESTAMP |  |
| db_uri | STRING |  |
| db_version | STRING |  |
| scanner_uri | STRING |  |
| scanner_version | STRING |  |
| origin | STRING |  |
| collector | STRING |  |
| document_ref | STRING |  |
| vulnerability_id | UUID |  |
| package_id | UUID |  |

#### Edges
| To | Edge Field | On Delete |
| --- | --- | --- |
| vulnerability_ids | certify_vulns_vulnerability_ids_vulnerability | vulnerability_id |
| package_versions | certify_vulns_package_versions_package | package_id |

#### Natural Keys (Composite Unique Constraints)
| Keys |
| --- |
| `[package, vulnerability, collector, scanner_uri, scanner_version, origin, db_uri, db_version]` |

### dependencies
#### Fields
| Field | Type | Comment |
| --- | --- | --- |
| id | UUID |  |
| dependency_type | STRING |  |
| justification | STRING |  |
| origin | STRING |  |
| collector | STRING |  |
| document_ref | STRING |  |
| package_id | UUID |  |
| dependent_package_version_id | UUID |  |

#### Edges
| To | Edge Field | On Delete |
| --- | --- | --- |
| package_versions | dependencies_package_versions_package | package_id |
| package_versions | dependencies_package_versions_dependent_package_version | dependent_package_version_id |

#### Natural Keys (Composite Unique Constraints)
| Keys |
| --- |
| `[dependency_type, justification, origin, collector, document_ref, package, dependent_package_version]` |

### has_metadata
#### Fields
| Field | Type | Comment |
| --- | --- | --- |
| id | UUID |  |
| timestamp | TIMESTAMP |  |
| key | STRING |  |
| value | STRING |  |
| justification | STRING |  |
| origin | STRING |  |
| collector | STRING |  |
| document_ref | STRING |  |
| source_id | UUID |  |
| package_version_id | UUID |  |
| package_name_id | UUID |  |
| artifact_id | UUID |  |

#### Edges
| To | Edge Field | On Delete |
| --- | --- | --- |
| source_names | has_metadata_source_names_source | source_id |
| package_versions | has_metadata_package_versions_package_version | package_version_id |
| package_names | has_metadata_package_names_all_versions | package_name_id |
| artifacts | has_metadata_artifacts_artifact | artifact_id |

#### Natural Keys (Composite Unique Constraints)
| Keys |
| --- |
| `[key, value, justification, origin, collector, timestamp, document_ref, source]` |
| `[key, value, justification, origin, collector, timestamp, document_ref, package_version]` |
| `[key, value, justification, origin, collector, timestamp, document_ref, package_name]` |
| `[key, value, justification, origin, collector, timestamp, document_ref, artifact]` |

### has_source_ats
#### Fields
| Field | Type | Comment |
| --- | --- | --- |
| id | UUID |  |
| known_since | TIMESTAMP |  |
| justification | STRING |  |
| origin | STRING |  |
| collector | STRING |  |
| document_ref | STRING |  |
| package_version_id | UUID |  |
| package_name_id | UUID |  |
| source_id | UUID |  |

#### Edges
| To | Edge Field | On Delete |
| --- | --- | --- |
| package_versions | has_source_ats_package_versions_package_version | package_version_id |
| package_names | has_source_ats_package_names_all_versions | package_name_id |
| source_names | has_source_ats_source_names_source | source_id |

#### Natural Keys (Composite Unique Constraints)
| Keys |
| --- |
| `[source, package_version, justification, origin, collector, known_since, document_ref]` |
| `[source, package_name, justification, origin, collector, known_since, document_ref]` |

### hash_equals
#### Fields
| Field | Type | Comment |
| --- | --- | --- |
| id | UUID |  |
| origin | STRING |  |
| collector | STRING |  |
| justification | STRING |  |
| document_ref | STRING |  |
| artifacts_hash | STRING | An opaque hash of the artifact IDs that are equal |
| art_id | UUID |  |
| equal_art_id | UUID |  |

#### Edges
| To | Edge Field | On Delete |
| --- | --- | --- |
| artifacts | hash_equals_artifacts_artifact_a | art_id |
| artifacts | hash_equals_artifacts_artifact_b | equal_art_id |

#### Natural Keys (Composite Unique Constraints)
| Keys |
| --- |
| `[art, equal_art, artifacts_hash, origin, justification, collector, document_ref]` |

### licenses
#### Fields
| Field | Type | Comment |
| --- | --- | --- |
| id | UUID |  |
| name | STRING |  |
| inline | STRING |  |
| list_version | STRING |  |
| inline_hash | STRING | An opaque hash on the linline text |
| list_version_hash | STRING | An opaque hash on the list_version text |

#### Natural Keys (Composite Unique Constraints)
| Keys |
| --- |
| `[name, inline_hash, list_version_hash]` |

### occurrences
#### Fields
| Field | Type | Comment |
| --- | --- | --- |
| id | UUID |  |
| justification | STRING | Justification for the attested relationship |
| origin | STRING | Document from which this attestation is generated from |
| collector | STRING | GUAC collector for the document |
| document_ref | STRING |  |
| artifact_id | UUID | The artifact in the relationship |
| package_id | UUID |  |
| source_id | UUID |  |

#### Edges
| To | Edge Field | On Delete |
| --- | --- | --- |
| artifacts | occurrences_artifacts_artifact | artifact_id |
| package_versions | occurrences_package_versions_package | package_id |
| source_names | occurrences_source_names_source | source_id |

#### Natural Keys (Composite Unique Constraints)
| Keys |
| --- |
| `[justification, origin, collector, document_ref, artifact, package]` |
| `[justification, origin, collector, document_ref, artifact, source]` |

### package_names
#### Fields
| Field | Type | Comment |
| --- | --- | --- |
| id | UUID |  |
| type | STRING | This node matches a pkg:<type> partial pURL |
| namespace | STRING | In the pURL representation, each PackageNamespace matches the pkg:<type>/<namespace>/ partial pURL |
| name | STRING |  |

#### Natural Keys (Composite Unique Constraints)
| Keys |
| --- |
| `[name, namespace, type]` |

### package_versions
#### Fields
| Field | Type | Comment |
| --- | --- | --- |
| id | UUID |  |
| version | STRING |  |
| subpath | STRING |  |
| qualifiers | STRING |  |
| hash | STRING | A SHA1 of the qualifiers, subpath, version fields after sorting keys, used to ensure uniqueness of version records. |
| name_id | UUID |  |

#### Edges
| To | Edge Field | On Delete |
| --- | --- | --- |
| package_names | package_versions_package_names_versions | name_id |

#### Natural Keys (Composite Unique Constraints)
| Keys |
| --- |
| `[hash, name]` |
| `[version, subpath, qualifiers, name]` |

### pkg_equals
#### Fields
| Field | Type | Comment |
| --- | --- | --- |
| id | UUID |  |
| origin | STRING |  |
| collector | STRING |  |
| document_ref | STRING |  |
| justification | STRING |  |
| packages_hash | STRING | An opaque hash of the package IDs that are equal |
| pkg_id | UUID |  |
| equal_pkg_id | UUID |  |

#### Edges
| To | Edge Field | On Delete |
| --- | --- | --- |
| package_versions | pkg_equals_package_versions_package_a | pkg_id |
| package_versions | pkg_equals_package_versions_package_b | equal_pkg_id |

#### Natural Keys (Composite Unique Constraints)
| Keys |
| --- |
| `[pkg, equal_pkg, packages_hash, origin, justification, collector, document_ref]` |

### point_of_contacts
#### Fields
| Field | Type | Comment |
| --- | --- | --- |
| id | UUID |  |
| email | STRING |  |
| info | STRING |  |
| since | TIMESTAMP |  |
| justification | STRING |  |
| origin | STRING |  |
| collector | STRING |  |
| document_ref | STRING |  |
| source_id | UUID |  |
| package_version_id | UUID |  |
| package_name_id | UUID |  |
| artifact_id | UUID |  |

#### Edges
| To | Edge Field | On Delete |
| --- | --- | --- |
| source_names | point_of_contacts_source_names_source | source_id |
| package_versions | point_of_contacts_package_versions_package_version | package_version_id |
| package_names | point_of_contacts_package_names_all_versions | package_name_id |
| artifacts | point_of_contacts_artifacts_artifact | artifact_id |

#### Natural Keys (Composite Unique Constraints)
| Keys |
| --- |
| `[since, email, info, justification, origin, collector, document_ref, source]` |
| `[since, email, info, justification, origin, collector, document_ref, package_version]` |
| `[since, email, info, justification, origin, collector, document_ref, package_name]` |
| `[since, email, info, justification, origin, collector, document_ref, artifact]` |

### slsa_attestations
#### Fields
| Field | Type | Comment |
| --- | --- | --- |
| id | UUID |  |
| build_type | STRING |  |
| slsa_predicate | STRING |  |
| slsa_version | STRING |  |
| started_on | TIMESTAMP |  |
| finished_on | TIMESTAMP |  |
| origin | STRING |  |
| collector | STRING |  |
| document_ref | STRING |  |
| built_from_hash | STRING |  |
| built_by_id | UUID |  |
| subject_id | UUID |  |

#### Edges
| To | Edge Field | On Delete |
| --- | --- | --- |
| builders | slsa_attestations_builders_built_by | built_by_id |
| artifacts | slsa_attestations_artifacts_subject | subject_id |
| artifacts | slsa_attestation_built_from |  |

#### Natural Keys (Composite Unique Constraints)
| Keys |
| --- |
| `[subject, origin, collector, document_ref, build_type, slsa_version, built_by, built_from_hash, started_on, finished_on]` |

### source_names
#### Fields
| Field | Type | Comment |
| --- | --- | --- |
| id | UUID |  |
| type | STRING |  |
| namespace | STRING |  |
| name | STRING |  |
| commit | STRING |  |
| tag | STRING |  |

#### Natural Keys (Composite Unique Constraints)
| Keys |
| --- |
| `[type, namespace, name, commit, tag]` |

### vuln_equals
#### Fields
| Field | Type | Comment |
| --- | --- | --- |
| id | UUID |  |
| justification | STRING |  |
| origin | STRING |  |
| collector | STRING |  |
| document_ref | STRING |  |
| vulnerabilities_hash | STRING | An opaque hash of the vulnerability IDs that are equal |
| vuln_id | UUID |  |
| equal_vuln_id | UUID |  |

#### Edges
| To | Edge Field | On Delete |
| --- | --- | --- |
| vulnerability_ids | vuln_equals_vulnerability_ids_vulnerability_a | vuln_id |
| vulnerability_ids | vuln_equals_vulnerability_ids_vulnerability_b | equal_vuln_id |

#### Natural Keys (Composite Unique Constraints)
| Keys |
| --- |
| `[vuln, equal_vuln, vulnerabilities_hash, justification, origin, collector, document_ref]` |

### vulnerability_ids
#### Fields
| Field | Type | Comment |
| --- | --- | --- |
| id | UUID |  |
| vulnerability_id | STRING |  |
| type | STRING |  |

#### Natural Keys (Composite Unique Constraints)
| Keys |
| --- |
| `[vulnerability_id, type]` |

### vulnerability_metadata
#### Fields
| Field | Type | Comment |
| --- | --- | --- |
| id | UUID |  |
| score_type | STRING |  |
| score_value | DOUBLE |  |
| timestamp | TIMESTAMP |  |
| origin | STRING |  |
| collector | STRING |  |
| document_ref | STRING |  |
| vulnerability_id_id | UUID |  |

#### Edges
| To | Edge Field | On Delete |
| --- | --- | --- |
| vulnerability_ids | vulnerability_metadata_vulnerability_ids_vulnerability_id | vulnerability_id_id |

#### Natural Keys (Composite Unique Constraints)
| Keys |
| --- |
| `[vulnerability_id, score_type, score_value, timestamp, origin, collector, document_ref]` |

