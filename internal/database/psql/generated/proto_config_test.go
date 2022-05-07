// Code generated by SQLBoiler 4.11.0 (https://github.com/volatiletech/sqlboiler). DO NOT EDIT.
// This file is meant to be re-generated in place and/or deleted at any time.

package models

import (
	"bytes"
	"context"
	"reflect"
	"testing"

	"github.com/volatiletech/randomize"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"github.com/volatiletech/strmangle"
)

var (
	// Relationships sometimes use the reflection helper queries.Equal/queries.Assign
	// so force a package dependency in case they don't.
	_ = queries.Equal
)

func testProtoConfigs(t *testing.T) {
	t.Parallel()

	query := ProtoConfigs()

	if query.Query == nil {
		t.Error("expected a query, got nothing")
	}
}

func testProtoConfigsDelete(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &ProtoConfig{}
	if err = randomize.Struct(seed, o, protoConfigDBTypes, true, protoConfigColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ProtoConfig struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if rowsAff, err := o.Delete(ctx, tx); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only have deleted one row, but affected:", rowsAff)
	}

	count, err := ProtoConfigs().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testProtoConfigsQueryDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &ProtoConfig{}
	if err = randomize.Struct(seed, o, protoConfigDBTypes, true, protoConfigColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ProtoConfig struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if rowsAff, err := ProtoConfigs().DeleteAll(ctx, tx); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only have deleted one row, but affected:", rowsAff)
	}

	count, err := ProtoConfigs().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testProtoConfigsSliceDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &ProtoConfig{}
	if err = randomize.Struct(seed, o, protoConfigDBTypes, true, protoConfigColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ProtoConfig struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice := ProtoConfigSlice{o}

	if rowsAff, err := slice.DeleteAll(ctx, tx); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only have deleted one row, but affected:", rowsAff)
	}

	count, err := ProtoConfigs().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testProtoConfigsExists(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &ProtoConfig{}
	if err = randomize.Struct(seed, o, protoConfigDBTypes, true, protoConfigColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ProtoConfig struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	e, err := ProtoConfigExists(ctx, tx, o.Pubkey)
	if err != nil {
		t.Errorf("Unable to check if ProtoConfig exists: %s", err)
	}
	if !e {
		t.Errorf("Expected ProtoConfigExists to return true, but got false.")
	}
}

func testProtoConfigsFind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &ProtoConfig{}
	if err = randomize.Struct(seed, o, protoConfigDBTypes, true, protoConfigColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ProtoConfig struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	protoConfigFound, err := FindProtoConfig(ctx, tx, o.Pubkey)
	if err != nil {
		t.Error(err)
	}

	if protoConfigFound == nil {
		t.Error("want a record, got nil")
	}
}

func testProtoConfigsBind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &ProtoConfig{}
	if err = randomize.Struct(seed, o, protoConfigDBTypes, true, protoConfigColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ProtoConfig struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if err = ProtoConfigs().Bind(ctx, tx, o); err != nil {
		t.Error(err)
	}
}

func testProtoConfigsOne(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &ProtoConfig{}
	if err = randomize.Struct(seed, o, protoConfigDBTypes, true, protoConfigColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ProtoConfig struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if x, err := ProtoConfigs().One(ctx, tx); err != nil {
		t.Error(err)
	} else if x == nil {
		t.Error("expected to get a non nil record")
	}
}

func testProtoConfigsAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	protoConfigOne := &ProtoConfig{}
	protoConfigTwo := &ProtoConfig{}
	if err = randomize.Struct(seed, protoConfigOne, protoConfigDBTypes, false, protoConfigColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ProtoConfig struct: %s", err)
	}
	if err = randomize.Struct(seed, protoConfigTwo, protoConfigDBTypes, false, protoConfigColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ProtoConfig struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = protoConfigOne.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}
	if err = protoConfigTwo.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice, err := ProtoConfigs().All(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 2 {
		t.Error("want 2 records, got:", len(slice))
	}
}

func testProtoConfigsCount(t *testing.T) {
	t.Parallel()

	var err error
	seed := randomize.NewSeed()
	protoConfigOne := &ProtoConfig{}
	protoConfigTwo := &ProtoConfig{}
	if err = randomize.Struct(seed, protoConfigOne, protoConfigDBTypes, false, protoConfigColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ProtoConfig struct: %s", err)
	}
	if err = randomize.Struct(seed, protoConfigTwo, protoConfigDBTypes, false, protoConfigColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ProtoConfig struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = protoConfigOne.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}
	if err = protoConfigTwo.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := ProtoConfigs().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 2 {
		t.Error("want 2 records, got:", count)
	}
}

func protoConfigBeforeInsertHook(ctx context.Context, e boil.ContextExecutor, o *ProtoConfig) error {
	*o = ProtoConfig{}
	return nil
}

func protoConfigAfterInsertHook(ctx context.Context, e boil.ContextExecutor, o *ProtoConfig) error {
	*o = ProtoConfig{}
	return nil
}

func protoConfigAfterSelectHook(ctx context.Context, e boil.ContextExecutor, o *ProtoConfig) error {
	*o = ProtoConfig{}
	return nil
}

func protoConfigBeforeUpdateHook(ctx context.Context, e boil.ContextExecutor, o *ProtoConfig) error {
	*o = ProtoConfig{}
	return nil
}

func protoConfigAfterUpdateHook(ctx context.Context, e boil.ContextExecutor, o *ProtoConfig) error {
	*o = ProtoConfig{}
	return nil
}

func protoConfigBeforeDeleteHook(ctx context.Context, e boil.ContextExecutor, o *ProtoConfig) error {
	*o = ProtoConfig{}
	return nil
}

func protoConfigAfterDeleteHook(ctx context.Context, e boil.ContextExecutor, o *ProtoConfig) error {
	*o = ProtoConfig{}
	return nil
}

func protoConfigBeforeUpsertHook(ctx context.Context, e boil.ContextExecutor, o *ProtoConfig) error {
	*o = ProtoConfig{}
	return nil
}

func protoConfigAfterUpsertHook(ctx context.Context, e boil.ContextExecutor, o *ProtoConfig) error {
	*o = ProtoConfig{}
	return nil
}

func testProtoConfigsHooks(t *testing.T) {
	t.Parallel()

	var err error

	ctx := context.Background()
	empty := &ProtoConfig{}
	o := &ProtoConfig{}

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, o, protoConfigDBTypes, false); err != nil {
		t.Errorf("Unable to randomize ProtoConfig object: %s", err)
	}

	AddProtoConfigHook(boil.BeforeInsertHook, protoConfigBeforeInsertHook)
	if err = o.doBeforeInsertHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doBeforeInsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeInsertHook function to empty object, but got: %#v", o)
	}
	protoConfigBeforeInsertHooks = []ProtoConfigHook{}

	AddProtoConfigHook(boil.AfterInsertHook, protoConfigAfterInsertHook)
	if err = o.doAfterInsertHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterInsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterInsertHook function to empty object, but got: %#v", o)
	}
	protoConfigAfterInsertHooks = []ProtoConfigHook{}

	AddProtoConfigHook(boil.AfterSelectHook, protoConfigAfterSelectHook)
	if err = o.doAfterSelectHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterSelectHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterSelectHook function to empty object, but got: %#v", o)
	}
	protoConfigAfterSelectHooks = []ProtoConfigHook{}

	AddProtoConfigHook(boil.BeforeUpdateHook, protoConfigBeforeUpdateHook)
	if err = o.doBeforeUpdateHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doBeforeUpdateHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeUpdateHook function to empty object, but got: %#v", o)
	}
	protoConfigBeforeUpdateHooks = []ProtoConfigHook{}

	AddProtoConfigHook(boil.AfterUpdateHook, protoConfigAfterUpdateHook)
	if err = o.doAfterUpdateHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterUpdateHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterUpdateHook function to empty object, but got: %#v", o)
	}
	protoConfigAfterUpdateHooks = []ProtoConfigHook{}

	AddProtoConfigHook(boil.BeforeDeleteHook, protoConfigBeforeDeleteHook)
	if err = o.doBeforeDeleteHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doBeforeDeleteHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeDeleteHook function to empty object, but got: %#v", o)
	}
	protoConfigBeforeDeleteHooks = []ProtoConfigHook{}

	AddProtoConfigHook(boil.AfterDeleteHook, protoConfigAfterDeleteHook)
	if err = o.doAfterDeleteHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterDeleteHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterDeleteHook function to empty object, but got: %#v", o)
	}
	protoConfigAfterDeleteHooks = []ProtoConfigHook{}

	AddProtoConfigHook(boil.BeforeUpsertHook, protoConfigBeforeUpsertHook)
	if err = o.doBeforeUpsertHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doBeforeUpsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeUpsertHook function to empty object, but got: %#v", o)
	}
	protoConfigBeforeUpsertHooks = []ProtoConfigHook{}

	AddProtoConfigHook(boil.AfterUpsertHook, protoConfigAfterUpsertHook)
	if err = o.doAfterUpsertHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterUpsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterUpsertHook function to empty object, but got: %#v", o)
	}
	protoConfigAfterUpsertHooks = []ProtoConfigHook{}
}

func testProtoConfigsInsert(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &ProtoConfig{}
	if err = randomize.Struct(seed, o, protoConfigDBTypes, true, protoConfigColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ProtoConfig struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := ProtoConfigs().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testProtoConfigsInsertWhitelist(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &ProtoConfig{}
	if err = randomize.Struct(seed, o, protoConfigDBTypes, true); err != nil {
		t.Errorf("Unable to randomize ProtoConfig struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Whitelist(protoConfigColumnsWithoutDefault...)); err != nil {
		t.Error(err)
	}

	count, err := ProtoConfigs().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testProtoConfigToManyVaults(t *testing.T) {
	var err error
	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()

	var a ProtoConfig
	var b, c Vault

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, protoConfigDBTypes, true, protoConfigColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ProtoConfig struct: %s", err)
	}

	if err := a.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	if err = randomize.Struct(seed, &b, vaultDBTypes, false, vaultColumnsWithDefault...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &c, vaultDBTypes, false, vaultColumnsWithDefault...); err != nil {
		t.Fatal(err)
	}

	b.ProtoConfig = a.Pubkey
	c.ProtoConfig = a.Pubkey

	if err = b.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}
	if err = c.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	check, err := a.Vaults().All(ctx, tx)
	if err != nil {
		t.Fatal(err)
	}

	bFound, cFound := false, false
	for _, v := range check {
		if v.ProtoConfig == b.ProtoConfig {
			bFound = true
		}
		if v.ProtoConfig == c.ProtoConfig {
			cFound = true
		}
	}

	if !bFound {
		t.Error("expected to find b")
	}
	if !cFound {
		t.Error("expected to find c")
	}

	slice := ProtoConfigSlice{&a}
	if err = a.L.LoadVaults(ctx, tx, false, (*[]*ProtoConfig)(&slice), nil); err != nil {
		t.Fatal(err)
	}
	if got := len(a.R.Vaults); got != 2 {
		t.Error("number of eager loaded records wrong, got:", got)
	}

	a.R.Vaults = nil
	if err = a.L.LoadVaults(ctx, tx, true, &a, nil); err != nil {
		t.Fatal(err)
	}
	if got := len(a.R.Vaults); got != 2 {
		t.Error("number of eager loaded records wrong, got:", got)
	}

	if t.Failed() {
		t.Logf("%#v", check)
	}
}

func testProtoConfigToManyAddOpVaults(t *testing.T) {
	var err error

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()

	var a ProtoConfig
	var b, c, d, e Vault

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, protoConfigDBTypes, false, strmangle.SetComplement(protoConfigPrimaryKeyColumns, protoConfigColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	foreigners := []*Vault{&b, &c, &d, &e}
	for _, x := range foreigners {
		if err = randomize.Struct(seed, x, vaultDBTypes, false, strmangle.SetComplement(vaultPrimaryKeyColumns, vaultColumnsWithoutDefault)...); err != nil {
			t.Fatal(err)
		}
	}

	if err := a.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}
	if err = b.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}
	if err = c.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	foreignersSplitByInsertion := [][]*Vault{
		{&b, &c},
		{&d, &e},
	}

	for i, x := range foreignersSplitByInsertion {
		err = a.AddVaults(ctx, tx, i != 0, x...)
		if err != nil {
			t.Fatal(err)
		}

		first := x[0]
		second := x[1]

		if a.Pubkey != first.ProtoConfig {
			t.Error("foreign key was wrong value", a.Pubkey, first.ProtoConfig)
		}
		if a.Pubkey != second.ProtoConfig {
			t.Error("foreign key was wrong value", a.Pubkey, second.ProtoConfig)
		}

		if first.R.VaultProtoConfig != &a {
			t.Error("relationship was not added properly to the foreign slice")
		}
		if second.R.VaultProtoConfig != &a {
			t.Error("relationship was not added properly to the foreign slice")
		}

		if a.R.Vaults[i*2] != first {
			t.Error("relationship struct slice not set to correct value")
		}
		if a.R.Vaults[i*2+1] != second {
			t.Error("relationship struct slice not set to correct value")
		}

		count, err := a.Vaults().Count(ctx, tx)
		if err != nil {
			t.Fatal(err)
		}
		if want := int64((i + 1) * 2); count != want {
			t.Error("want", want, "got", count)
		}
	}
}

func testProtoConfigsReload(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &ProtoConfig{}
	if err = randomize.Struct(seed, o, protoConfigDBTypes, true, protoConfigColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ProtoConfig struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if err = o.Reload(ctx, tx); err != nil {
		t.Error(err)
	}
}

func testProtoConfigsReloadAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &ProtoConfig{}
	if err = randomize.Struct(seed, o, protoConfigDBTypes, true, protoConfigColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ProtoConfig struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice := ProtoConfigSlice{o}

	if err = slice.ReloadAll(ctx, tx); err != nil {
		t.Error(err)
	}
}

func testProtoConfigsSelect(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &ProtoConfig{}
	if err = randomize.Struct(seed, o, protoConfigDBTypes, true, protoConfigColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ProtoConfig struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice, err := ProtoConfigs().All(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 1 {
		t.Error("want one record, got:", len(slice))
	}
}

var (
	protoConfigDBTypes = map[string]string{`Pubkey`: `character varying`, `Granularity`: `numeric`, `TriggerDcaSpread`: `smallint`, `BaseWithdrawalSpread`: `smallint`}
	_                  = bytes.MinRead
)

func testProtoConfigsUpdate(t *testing.T) {
	t.Parallel()

	if 0 == len(protoConfigPrimaryKeyColumns) {
		t.Skip("Skipping table with no primary key columns")
	}
	if len(protoConfigAllColumns) == len(protoConfigPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	o := &ProtoConfig{}
	if err = randomize.Struct(seed, o, protoConfigDBTypes, true, protoConfigColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ProtoConfig struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := ProtoConfigs().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, o, protoConfigDBTypes, true, protoConfigPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize ProtoConfig struct: %s", err)
	}

	if rowsAff, err := o.Update(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only affect one row but affected", rowsAff)
	}
}

func testProtoConfigsSliceUpdateAll(t *testing.T) {
	t.Parallel()

	if len(protoConfigAllColumns) == len(protoConfigPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	o := &ProtoConfig{}
	if err = randomize.Struct(seed, o, protoConfigDBTypes, true, protoConfigColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ProtoConfig struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := ProtoConfigs().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, o, protoConfigDBTypes, true, protoConfigPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize ProtoConfig struct: %s", err)
	}

	// Remove Primary keys and unique columns from what we plan to update
	var fields []string
	if strmangle.StringSliceMatch(protoConfigAllColumns, protoConfigPrimaryKeyColumns) {
		fields = protoConfigAllColumns
	} else {
		fields = strmangle.SetComplement(
			protoConfigAllColumns,
			protoConfigPrimaryKeyColumns,
		)
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	typ := reflect.TypeOf(o).Elem()
	n := typ.NumField()

	updateMap := M{}
	for _, col := range fields {
		for i := 0; i < n; i++ {
			f := typ.Field(i)
			if f.Tag.Get("boil") == col {
				updateMap[col] = value.Field(i).Interface()
			}
		}
	}

	slice := ProtoConfigSlice{o}
	if rowsAff, err := slice.UpdateAll(ctx, tx, updateMap); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("wanted one record updated but got", rowsAff)
	}
}

func testProtoConfigsUpsert(t *testing.T) {
	t.Parallel()

	if len(protoConfigAllColumns) == len(protoConfigPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	// Attempt the INSERT side of an UPSERT
	o := ProtoConfig{}
	if err = randomize.Struct(seed, &o, protoConfigDBTypes, true); err != nil {
		t.Errorf("Unable to randomize ProtoConfig struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Upsert(ctx, tx, false, nil, boil.Infer(), boil.Infer()); err != nil {
		t.Errorf("Unable to upsert ProtoConfig: %s", err)
	}

	count, err := ProtoConfigs().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}

	// Attempt the UPDATE side of an UPSERT
	if err = randomize.Struct(seed, &o, protoConfigDBTypes, false, protoConfigPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize ProtoConfig struct: %s", err)
	}

	if err = o.Upsert(ctx, tx, true, nil, boil.Infer(), boil.Infer()); err != nil {
		t.Errorf("Unable to upsert ProtoConfig: %s", err)
	}

	count, err = ProtoConfigs().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}
}
