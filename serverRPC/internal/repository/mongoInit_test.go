package repository

/*
	func TestInsert(t *testing.T) {
		db := mongomock.NewDB()

		collection := db.Collection("someCollection")

		type mockBehavior func(ctx context.Context, products []domain.Product)

		testTables := []struct {
			name              string
			mockBehavior      mockBehavior
			products          []domain.Product
			productsInterface []interface{}
		}{
			{
				name: "Valid",
				mockBehavior: func(ctx context.Context, products []domain.Product) {
					collection.Insert()
				},
			},
		}
	}

	func TestMongoList(t *testing.T) {
		db := mongomock.NewDB()

		collection := db.Collection("someCollection")

		type mockBehavior func(ctx context.Context, sort domain.SortParams, products []domain.Product)

		testTables := []struct {
			name              string
			mockBehavior      mockBehavior
			products          []domain.Product
			productsInterface []interface{}
		}{
			{
				name: "Valid",
				mockBehavior: func(ctx context.Context, sort domain.SortParams, products []domain.Product) {
					opts := options.Find()
					sortOpts := bson.M{sort.SortField: sort.SortAsc}

					opts.SetSort(sortOpts)
					opts.SetSkip(int64(sort.PagingOffset))
					opts.SetLimit(int64(sort.PagingLimit))

					collection.Find(ctx, bson.M{})
				},
			},
		}
	}

func TestMongoGetByName(t *testing.T) {
	db := mongomock.NewDB()

	collection := db.Collection("someCollection")

	type mockBehavior func(ctx context.Context, product domain.Product)

	testTables := []struct {
		name         string
		mockBehavior mockBehavior
		products     domain.Product
		isErr        bool
	}{
		{
			name: "Valid",
			mockBehavior: func(ctx context.Context, product domain.Product) {
				filter := bson.M{"name": product.Name}

				err := collection.FindFirst(ctx, filter)
				if err != nil {
					assert.Error(t, err)
				}
			},
			products: domain.Product{
				Id:   1,
				Name: "name",
				Price: func() primitive.Decimal128 {
					got, _ := primitive.ParseDecimal128("50.00")
					return got
				}(),
			},
		},
	}

	logger := logger.GetLogger()

	product := &MongoBackend{
		db:     db,
		logger: logger,
		dbname: "assets",
	}

	repo := Repository{
		Product: product,
	}

	for _, table := range testTables {
		t.Run(table.name, func(t *testing.T) {
			table.mockBehavior()

			//err :=
			if table.isErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
*/
