package types

import (
	"log"
	"testing"

	"github.com/VKCOM/noverify/src/linttest/assert"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

type TestCase struct {
	T1, T2 string
	Result CompatibleResult
	Ok     bool
}

func TestCompatible(t *testing.T) {
	cases := []TestCase{
		{
			T1: "", T2: "int",
			Ok: true,
		},
		{
			T1: "int", T2: "",
			Ok: true,
		},
		// int and alias.
		{
			T1: "int", T2: "int",
			Ok: true,
		},
		{
			T1: "integer", T2: "int",
			Ok: true,
		},
		{
			T1: "long", T2: "int",
			Ok: true,
		},
		{
			T1: "long", T2: "integer",
			Ok: true,
		},
		{
			T1: "mixed", T2: "int",
			Ok: true,
		},
		{
			T1: "string", T2: "int",
			Ok: false,
		},

		// int and float.
		{
			T1: "int", T2: "float",
			Ok: false,
			Result: CompatibleResult{
				IntFloat: true,
			},
		},
		{
			T1: "float", T2: "int",
			Ok: false,
			Result: CompatibleResult{
				FloatInt: true,
			},
		},

		// bool and alias.
		{
			T1: "bool", T2: "bool",
			Ok: true,
		},
		{
			T1: "boolean", T2: "bool",
			Ok: true,
		},
		// bool and true/false.
		{
			T1: "bool", T2: "true",
			Ok: true,
		},
		{
			T1: "true", T2: "bool",
			Ok: true,
		},
		{
			T1: "bool", T2: "false",
			Ok: false,
			Result: CompatibleResult{
				BoolFalse: true,
			},
		},
		{
			T1: "false", T2: "bool",
			Ok: false,
			Result: CompatibleResult{
				FalseBool: true,
			},
		},
		{
			T1: "mixed", T2: "bool",
			Ok: true,
		},
		{
			T1: "string", T2: "bool",
			Ok: false,
		},

		// float and alias.
		{
			T1: "float", T2: "float",
			Ok: true,
		},
		{
			T1: "real", T2: "float",
			Ok: true,
		},
		{
			T1: "double", T2: "float",
			Ok: true,
		},
		{
			T1: "double", T2: "real",
			Ok: true,
		},
		{
			T1: "mixed", T2: "float",
			Ok: true,
		},
		{
			T1: "string", T2: "float",
			Ok: false,
		},

		// string.
		{
			T1: "string", T2: "string",
			Ok: true,
		},
		{
			T1: "mixed", T2: "string",
			Ok: true,
		},

		// Arrays.
		{
			T1: "int[]",
			T2: "int[]",
			Ok: true,
		},
		{
			T1: "int[][]",
			T2: "int[][]",
			Ok: true,
		},
		{
			T1: "double[][]",
			T2: "float[][]",
			Ok: true,
		},
		{
			T1: "int[]",
			T2: "float[]",
			Ok: false,
			Result: CompatibleResult{
				ArraysTypeMismatch: true,
				ArrayCheckResult: &CompatibleResult{
					IntFloat: true,
				},
			},
		},
		{
			T1: "bool[]",
			T2: "false[]",
			Ok: false,
			Result: CompatibleResult{
				ArraysTypeMismatch: true,
				ArrayCheckResult: &CompatibleResult{
					BoolFalse: true,
				},
			},
		},
		{
			T1: "bool[]",
			T2: "string",
			Ok: false,
			Result: CompatibleResult{
				ArrayAndType: true,
			},
		},
		{
			T1: "string",
			T2: "bool[]",
			Ok: false,
			Result: CompatibleResult{
				TypeAndArray: true,
			},
		},
		{
			T1: "float",
			T2: "mixed[]",
			Ok: false,
			Result: CompatibleResult{
				TypeAndArray: true,
			},
		},
		{
			T1: "mixed[]",
			T2: "string",
			Ok: false,
			Result: CompatibleResult{
				ArrayAndType: true,
			},
		},
		{
			T1: "int[]",
			T2: "mixed[]",
			Ok: true,
		},
		{
			T1: "int[][]",
			T2: "mixed[]",
			Ok: true,
		},
		{
			T1: "int[][]",
			T2: "mixed",
			Ok: true,
		},

		// Classes.
		{
			T1: `\SimpleClass`,
			T2: `\SimpleClass`,
			Ok: true,
		},
		{
			T1: `\SimpleClass`,
			T2: `object`,
			Ok: true,
		},
		{
			T1: `object`,
			T2: `\SimpleClass`,
			Ok: true,
		},
		{
			T1: `\SimpleClassWithSimpleIface`,
			T2: `\SimpleIface`,
			Ok: true,
			Result: CompatibleResult{
				ClassAndInterface: true,
			},
		},
		{
			T1: `\SimpleIface`,
			T2: `\SimpleClassWithSimpleIface`,
			Ok: true,
			Result: CompatibleResult{
				InterfaceAndClass: true,
			},
		},
		{
			T1: `\SimpleClassWithSimpleIface`,
			T2: `\SimpleIface2`, // Not implements.
			Ok: false,
			Result: CompatibleResult{
				ClassAndInterface: true,
			},
		},

		{
			T1: `\DerivedClassWithClassWithSimpleIface`,
			T2: `\SimpleIface`,
			Ok: true,
			Result: CompatibleResult{
				ClassAndInterface: true,
			},
		},
		{
			T1: `\SimpleIface`,
			T2: `\DerivedClassWithClassWithSimpleIface`,
			Ok: true,
			Result: CompatibleResult{
				InterfaceAndClass: true,
			},
		},
		{
			T1: `\DerivedClassWithClassWithSimpleIface`,
			T2: `\SimpleIface2`, // Not implements.
			Ok: false,
			Result: CompatibleResult{
				ClassAndInterface: true,
			},
		},
		{
			T1: `\DerivedClassWithClassWithDerivedIface`,
			T2: `\DerivedIface`,
			Ok: true,
			Result: CompatibleResult{
				ClassAndInterface: true,
			},
		},
		{
			T1: `\DerivedClassWithClassWithDerivedIface`,
			T2: `\BaseIface`,
			Ok: true,
			Result: CompatibleResult{
				ClassAndInterface: true,
			},
		},

		// Extends.
		{
			T1: `\BaseClass`,
			T2: `\DerivedClassFromBaseClass`,
			Ok: true,
			Result: CompatibleResult{
				ParentAndClass: true,
			},
		},
		{
			T1: `\DerivedClassFromBaseClass`,
			T2: `\BaseClass`,
			Ok: true,
			Result: CompatibleResult{
				ClassAndParent: true,
			},
		},
		{
			T1: `\DerivedClassFromBaseClass`,
			T2: `\BaseClass2`,
			Ok: false, // Not extends.
		},

		// With not class.
		{
			T1: `\DerivedClassFromBaseClass`,
			T2: `int`,
			Ok: false,
			Result: CompatibleResult{
				ClassAndNotClass: true,
			},
		},
		{
			T1: `int`,
			T2: `\DerivedClassFromBaseClass`,
			Ok: false,
			Result: CompatibleResult{
				NotClassAndClass: true,
			},
		},
	}

	classes := []ClassData{
		{
			Name: `\SimpleClass`,
		},
		{
			Name:        `\SimpleIface`,
			IsInterface: true,
		},
		{
			Name:        `\SimpleIface2`,
			IsInterface: true,
		},
		{
			Name:       `\SimpleClassWithSimpleIface`,
			Interfaces: map[string]struct{}{`\SimpleIface`: {}},
		},

		{
			Name:       `\BaseClassWithIface`,
			Interfaces: map[string]struct{}{`\SimpleIface`: {}},
		},
		{
			Name:   `\DerivedClassWithClassWithSimpleIface`,
			Parent: `\BaseClassWithIface`,
		},

		{
			Name: `\BaseClass`,
		},

		{
			Name: `\BaseClass2`,
		},

		{
			Name:   `\DerivedClassFromBaseClass`,
			Parent: `\BaseClass`,
		},

		{
			Name:        `\BaseIface`,
			IsInterface: true,
		},
		{
			Name:        `\DerivedIface`,
			IsInterface: true,
			Interfaces:  map[string]struct{}{`\BaseIface`: {}},
		},
		{
			Name:       `\BaseClassWithDerivedIface`,
			Interfaces: map[string]struct{}{`\DerivedIface`: {}},
		},
		{
			Name:   `\DerivedClassWithClassWithDerivedIface`,
			Parent: `\BaseClassWithDerivedIface`,
		},
	}

	comparator := Compatible{
		ClassDataProvider: func(name string) (ClassData, bool) {
			for _, class := range classes {
				if class.Name == name {
					return class, true
				}
			}

			return ClassData{}, false
		},
	}

	for i, testCase := range cases {
		map1 := NewMap(testCase.T1)
		map2 := NewMap(testCase.T2)

		testCase.Result.IsCompatible = testCase.Ok

		result := comparator.CompatibleTypes(map1, map2)
		if !cmp.Equal(result, testCase.Result, cmpopts.IgnoreTypes(NewMap(""))) {
			log.Printf("Case #%d:\nT1: %s\nT2: %s\n", i, map1, map2)
			assert.DeepEqual(t, result, testCase.Result, cmpopts.IgnoreTypes(NewMap("")))
		}
	}
}