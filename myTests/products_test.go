package myTests

import "testing"

type testFindOneProduct struct {
	ProductId string
	isError   bool
	expected  string
}

func TestFindOneProduct(t *testing.T) {
	tests := []testFindOneProduct{
		{
			ProductId: "P0000999",
			isError:   true,
			expected:  "get product failed: sql: no rows in result set",
		},
		{
			ProductId: "P000001",
			isError:   false,
			expected:  `{"id":"P000001","title":"Coffee","description":"Just a food \u0026 beverage product","category":{"id":1,"title":"food \u0026 beverage"},"created_at":"2023-11-15T22:21:05.247324","updated_at":"2023-11-15T22:21:05.247324","price":150,"images":[{"id":"c580fe73-afb3-47d1-a9df-eed24fdaea9b","filename":"fb1_1.jpg","url":"https://i.pinimg.com/564x/4a/1c/4a/4a1c4a9755e4d3bdfcb45a1c3a58712f.jpg"},{"id":"43bcd3fa-6f7f-4251-b196-f30ad4ea625e","filename":"fb1_2.jpg","url":"https://i.pinimg.com/564x/4a/1c/4a/4a1c4a9755e4d3bdfcb45a1c3a58712f.jpg"},{"id":"77d9e690-b722-4039-b0fe-5f7d9af0e6b4","filename":"fb1_3.jpg","url":"https://i.pinimg.com/564x/4a/1c/4a/4a1c4a9755e4d3bdfcb45a1c3a58712f.jpg"}]}`,
		},
	}

	productModule := SetupTest().ProductsModule()
	for _, test := range tests {
		if test.isError {
			if _, err := productModule.Usecase().FindOneProduct(test.ProductId); err.Error() != test.expected {
				t.Errorf("expected: %v, got: %v", test.expected, err.Error())
			}
		} else {
			result, err := productModule.Usecase().FindOneProduct(test.ProductId)
			if err != nil {
				t.Errorf("expected: %v, got: %v", nil, err.Error())
			}
			if CompressToJSON(result) != test.expected {
				t.Errorf("expected: %v, got: %v", test.expected, CompressToJSON(result))
			}
		}

	}
}
