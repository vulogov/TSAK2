package dsl

import (
	"fmt"

	. "github.com/glycerine/zygomys/zygo"
	"syreclabs.com/go/faker"
)

func FakeGen(env *Zlisp, name string, args []Sexp) (Sexp, error) {
	switch name {
	case "fake.Name":
		return &SexpStr{S: faker.Name().String()}, nil
	case "fake.Address":
		return &SexpStr{S: faker.Address().String()}, nil
	case "fake.Company":
		return &SexpStr{S: faker.Company().Name()}, nil
	case "fake.Phone":
		return &SexpStr{S: faker.PhoneNumber().PhoneNumber()}, nil
	}
	return SexpNull, fmt.Errorf("Do not know how to generate that request for fake data")
}

func FakeFunctions() map[string]ZlispUserFunction {
	return map[string]ZlispUserFunction{
		"fakename":        FakeGen,
		"fakeaddress":     FakeGen,
		"fakecompany":     FakeGen,
		"fakephonenumber": FakeGen,
	}
}

func FakePackageSetup(cfg *ZlispConfig, env *Zlisp) {
	myPkg := `(def fake (package "fake"
     { Name := fakename;
       Address := fakeaddress;
			 Company := fakecompany;
	     Phone := fakephonenumber;
     }
  ))`
	_, err := env.EvalString(myPkg)
	PanicOn(err)
}
