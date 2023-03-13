update "sites"
set
  "name" = :name,
  "origins" = :origins,
  "use_arweave" = :use_arweave,
  "updated_at" = NOW()
where
  "id" = :id
;