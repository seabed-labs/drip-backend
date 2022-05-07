// Code generated by SQLBoiler 4.11.0 (https://github.com/volatiletech/sqlboiler). DO NOT EDIT.
// This file is meant to be re-generated in place and/or deleted at any time.

package models

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/friendsofgo/errors"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"github.com/volatiletech/sqlboiler/v4/queries/qmhelper"
	"github.com/volatiletech/strmangle"
)

// UserPosition is an object representing the database table.
type UserPosition struct {
	Pubkey string `boil:"pubkey" json:"pubkey" toml:"pubkey" yaml:"pubkey"`
	Mint   string `boil:"mint" json:"mint" toml:"mint" yaml:"mint"`
	Amount bool   `boil:"amount" json:"amount" toml:"amount" yaml:"amount"`

	R *userPositionR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L userPositionL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var UserPositionColumns = struct {
	Pubkey string
	Mint   string
	Amount string
}{
	Pubkey: "pubkey",
	Mint:   "mint",
	Amount: "amount",
}

var UserPositionTableColumns = struct {
	Pubkey string
	Mint   string
	Amount string
}{
	Pubkey: "user_position.pubkey",
	Mint:   "user_position.mint",
	Amount: "user_position.amount",
}

// Generated where

var UserPositionWhere = struct {
	Pubkey whereHelperstring
	Mint   whereHelperstring
	Amount whereHelperbool
}{
	Pubkey: whereHelperstring{field: "\"user_position\".\"pubkey\""},
	Mint:   whereHelperstring{field: "\"user_position\".\"mint\""},
	Amount: whereHelperbool{field: "\"user_position\".\"amount\""},
}

// UserPositionRels is where relationship names are stored.
var UserPositionRels = struct {
}{}

// userPositionR is where relationships are stored.
type userPositionR struct {
}

// NewStruct creates a new relationship struct
func (*userPositionR) NewStruct() *userPositionR {
	return &userPositionR{}
}

// userPositionL is where Load methods for each relationship are stored.
type userPositionL struct{}

var (
	userPositionAllColumns            = []string{"pubkey", "mint", "amount"}
	userPositionColumnsWithoutDefault = []string{"pubkey", "mint", "amount"}
	userPositionColumnsWithDefault    = []string{}
	userPositionPrimaryKeyColumns     = []string{"pubkey"}
	userPositionGeneratedColumns      = []string{}
)

type (
	// UserPositionSlice is an alias for a slice of pointers to UserPosition.
	// This should almost always be used instead of []UserPosition.
	UserPositionSlice []*UserPosition
	// UserPositionHook is the signature for custom UserPosition hook methods
	UserPositionHook func(context.Context, boil.ContextExecutor, *UserPosition) error

	userPositionQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	userPositionType                 = reflect.TypeOf(&UserPosition{})
	userPositionMapping              = queries.MakeStructMapping(userPositionType)
	userPositionPrimaryKeyMapping, _ = queries.BindMapping(userPositionType, userPositionMapping, userPositionPrimaryKeyColumns)
	userPositionInsertCacheMut       sync.RWMutex
	userPositionInsertCache          = make(map[string]insertCache)
	userPositionUpdateCacheMut       sync.RWMutex
	userPositionUpdateCache          = make(map[string]updateCache)
	userPositionUpsertCacheMut       sync.RWMutex
	userPositionUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

var userPositionAfterSelectHooks []UserPositionHook

var userPositionBeforeInsertHooks []UserPositionHook
var userPositionAfterInsertHooks []UserPositionHook

var userPositionBeforeUpdateHooks []UserPositionHook
var userPositionAfterUpdateHooks []UserPositionHook

var userPositionBeforeDeleteHooks []UserPositionHook
var userPositionAfterDeleteHooks []UserPositionHook

var userPositionBeforeUpsertHooks []UserPositionHook
var userPositionAfterUpsertHooks []UserPositionHook

// doAfterSelectHooks executes all "after Select" hooks.
func (o *UserPosition) doAfterSelectHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range userPositionAfterSelectHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeInsertHooks executes all "before insert" hooks.
func (o *UserPosition) doBeforeInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range userPositionBeforeInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterInsertHooks executes all "after Insert" hooks.
func (o *UserPosition) doAfterInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range userPositionAfterInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpdateHooks executes all "before Update" hooks.
func (o *UserPosition) doBeforeUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range userPositionBeforeUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpdateHooks executes all "after Update" hooks.
func (o *UserPosition) doAfterUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range userPositionAfterUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeDeleteHooks executes all "before Delete" hooks.
func (o *UserPosition) doBeforeDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range userPositionBeforeDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterDeleteHooks executes all "after Delete" hooks.
func (o *UserPosition) doAfterDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range userPositionAfterDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpsertHooks executes all "before Upsert" hooks.
func (o *UserPosition) doBeforeUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range userPositionBeforeUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpsertHooks executes all "after Upsert" hooks.
func (o *UserPosition) doAfterUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range userPositionAfterUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// AddUserPositionHook registers your hook function for all future operations.
func AddUserPositionHook(hookPoint boil.HookPoint, userPositionHook UserPositionHook) {
	switch hookPoint {
	case boil.AfterSelectHook:
		userPositionAfterSelectHooks = append(userPositionAfterSelectHooks, userPositionHook)
	case boil.BeforeInsertHook:
		userPositionBeforeInsertHooks = append(userPositionBeforeInsertHooks, userPositionHook)
	case boil.AfterInsertHook:
		userPositionAfterInsertHooks = append(userPositionAfterInsertHooks, userPositionHook)
	case boil.BeforeUpdateHook:
		userPositionBeforeUpdateHooks = append(userPositionBeforeUpdateHooks, userPositionHook)
	case boil.AfterUpdateHook:
		userPositionAfterUpdateHooks = append(userPositionAfterUpdateHooks, userPositionHook)
	case boil.BeforeDeleteHook:
		userPositionBeforeDeleteHooks = append(userPositionBeforeDeleteHooks, userPositionHook)
	case boil.AfterDeleteHook:
		userPositionAfterDeleteHooks = append(userPositionAfterDeleteHooks, userPositionHook)
	case boil.BeforeUpsertHook:
		userPositionBeforeUpsertHooks = append(userPositionBeforeUpsertHooks, userPositionHook)
	case boil.AfterUpsertHook:
		userPositionAfterUpsertHooks = append(userPositionAfterUpsertHooks, userPositionHook)
	}
}

// One returns a single userPosition record from the query.
func (q userPositionQuery) One(ctx context.Context, exec boil.ContextExecutor) (*UserPosition, error) {
	o := &UserPosition{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(ctx, exec, o)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: failed to execute a one query for user_position")
	}

	if err := o.doAfterSelectHooks(ctx, exec); err != nil {
		return o, err
	}

	return o, nil
}

// All returns all UserPosition records from the query.
func (q userPositionQuery) All(ctx context.Context, exec boil.ContextExecutor) (UserPositionSlice, error) {
	var o []*UserPosition

	err := q.Bind(ctx, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "models: failed to assign all query results to UserPosition slice")
	}

	if len(userPositionAfterSelectHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterSelectHooks(ctx, exec); err != nil {
				return o, err
			}
		}
	}

	return o, nil
}

// Count returns the count of all UserPosition records in the query.
func (q userPositionQuery) Count(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to count user_position rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table.
func (q userPositionQuery) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "models: failed to check if user_position exists")
	}

	return count > 0, nil
}

// UserPositions retrieves all the records using an executor.
func UserPositions(mods ...qm.QueryMod) userPositionQuery {
	mods = append(mods, qm.From("\"user_position\""))
	q := NewQuery(mods...)
	if len(queries.GetSelect(q)) == 0 {
		queries.SetSelect(q, []string{"\"user_position\".*"})
	}

	return userPositionQuery{q}
}

// FindUserPosition retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindUserPosition(ctx context.Context, exec boil.ContextExecutor, pubkey string, selectCols ...string) (*UserPosition, error) {
	userPositionObj := &UserPosition{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"user_position\" where \"pubkey\"=$1", sel,
	)

	q := queries.Raw(query, pubkey)

	err := q.Bind(ctx, exec, userPositionObj)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: unable to select from user_position")
	}

	if err = userPositionObj.doAfterSelectHooks(ctx, exec); err != nil {
		return userPositionObj, err
	}

	return userPositionObj, nil
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *UserPosition) Insert(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error {
	if o == nil {
		return errors.New("models: no user_position provided for insertion")
	}

	var err error

	if err := o.doBeforeInsertHooks(ctx, exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(userPositionColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	userPositionInsertCacheMut.RLock()
	cache, cached := userPositionInsertCache[key]
	userPositionInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			userPositionAllColumns,
			userPositionColumnsWithDefault,
			userPositionColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(userPositionType, userPositionMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(userPositionType, userPositionMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO \"user_position\" (\"%s\") %%sVALUES (%s)%%s", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO \"user_position\" %sDEFAULT VALUES%s"
		}

		var queryOutput, queryReturning string

		if len(cache.retMapping) != 0 {
			queryReturning = fmt.Sprintf(" RETURNING \"%s\"", strings.Join(returnColumns, "\",\""))
		}

		cache.query = fmt.Sprintf(cache.query, queryOutput, queryReturning)
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	vals := queries.ValuesFromMapping(value, cache.valueMapping)

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.query)
		fmt.Fprintln(writer, vals)
	}

	if len(cache.retMapping) != 0 {
		err = exec.QueryRowContext(ctx, cache.query, vals...).Scan(queries.PtrsFromMapping(value, cache.retMapping)...)
	} else {
		_, err = exec.ExecContext(ctx, cache.query, vals...)
	}

	if err != nil {
		return errors.Wrap(err, "models: unable to insert into user_position")
	}

	if !cached {
		userPositionInsertCacheMut.Lock()
		userPositionInsertCache[key] = cache
		userPositionInsertCacheMut.Unlock()
	}

	return o.doAfterInsertHooks(ctx, exec)
}

// Update uses an executor to update the UserPosition.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *UserPosition) Update(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) (int64, error) {
	var err error
	if err = o.doBeforeUpdateHooks(ctx, exec); err != nil {
		return 0, err
	}
	key := makeCacheKey(columns, nil)
	userPositionUpdateCacheMut.RLock()
	cache, cached := userPositionUpdateCache[key]
	userPositionUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			userPositionAllColumns,
			userPositionPrimaryKeyColumns,
		)

		if !columns.IsWhitelist() {
			wl = strmangle.SetComplement(wl, []string{"created_at"})
		}
		if len(wl) == 0 {
			return 0, errors.New("models: unable to update user_position, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"user_position\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, wl),
			strmangle.WhereClause("\"", "\"", len(wl)+1, userPositionPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(userPositionType, userPositionMapping, append(wl, userPositionPrimaryKeyColumns...))
		if err != nil {
			return 0, err
		}
	}

	values := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), cache.valueMapping)

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.query)
		fmt.Fprintln(writer, values)
	}
	var result sql.Result
	result, err = exec.ExecContext(ctx, cache.query, values...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update user_position row")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by update for user_position")
	}

	if !cached {
		userPositionUpdateCacheMut.Lock()
		userPositionUpdateCache[key] = cache
		userPositionUpdateCacheMut.Unlock()
	}

	return rowsAff, o.doAfterUpdateHooks(ctx, exec)
}

// UpdateAll updates all rows with the specified column values.
func (q userPositionQuery) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	queries.SetUpdate(q.Query, cols)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all for user_position")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected for user_position")
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o UserPositionSlice) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	ln := int64(len(o))
	if ln == 0 {
		return 0, nil
	}

	if len(cols) == 0 {
		return 0, errors.New("models: update all requires at least one column argument")
	}

	colNames := make([]string, len(cols))
	args := make([]interface{}, len(cols))

	i := 0
	for name, value := range cols {
		colNames[i] = name
		args[i] = value
		i++
	}

	// Append all of the primary key values for each column
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), userPositionPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE \"user_position\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), len(colNames)+1, userPositionPrimaryKeyColumns, len(o)))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all in userPosition slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected all in update all userPosition")
	}
	return rowsAff, nil
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *UserPosition) Upsert(ctx context.Context, exec boil.ContextExecutor, updateOnConflict bool, conflictColumns []string, updateColumns, insertColumns boil.Columns) error {
	if o == nil {
		return errors.New("models: no user_position provided for upsert")
	}

	if err := o.doBeforeUpsertHooks(ctx, exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(userPositionColumnsWithDefault, o)

	// Build cache key in-line uglily - mysql vs psql problems
	buf := strmangle.GetBuffer()
	if updateOnConflict {
		buf.WriteByte('t')
	} else {
		buf.WriteByte('f')
	}
	buf.WriteByte('.')
	for _, c := range conflictColumns {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	buf.WriteString(strconv.Itoa(updateColumns.Kind))
	for _, c := range updateColumns.Cols {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	buf.WriteString(strconv.Itoa(insertColumns.Kind))
	for _, c := range insertColumns.Cols {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	for _, c := range nzDefaults {
		buf.WriteString(c)
	}
	key := buf.String()
	strmangle.PutBuffer(buf)

	userPositionUpsertCacheMut.RLock()
	cache, cached := userPositionUpsertCache[key]
	userPositionUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, ret := insertColumns.InsertColumnSet(
			userPositionAllColumns,
			userPositionColumnsWithDefault,
			userPositionColumnsWithoutDefault,
			nzDefaults,
		)

		update := updateColumns.UpdateColumnSet(
			userPositionAllColumns,
			userPositionPrimaryKeyColumns,
		)

		if updateOnConflict && len(update) == 0 {
			return errors.New("models: unable to upsert user_position, could not build update column list")
		}

		conflict := conflictColumns
		if len(conflict) == 0 {
			conflict = make([]string, len(userPositionPrimaryKeyColumns))
			copy(conflict, userPositionPrimaryKeyColumns)
		}
		cache.query = buildUpsertQueryPostgres(dialect, "\"user_position\"", updateOnConflict, ret, update, conflict, insert)

		cache.valueMapping, err = queries.BindMapping(userPositionType, userPositionMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(userPositionType, userPositionMapping, ret)
			if err != nil {
				return err
			}
		}
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	vals := queries.ValuesFromMapping(value, cache.valueMapping)
	var returns []interface{}
	if len(cache.retMapping) != 0 {
		returns = queries.PtrsFromMapping(value, cache.retMapping)
	}

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.query)
		fmt.Fprintln(writer, vals)
	}
	if len(cache.retMapping) != 0 {
		err = exec.QueryRowContext(ctx, cache.query, vals...).Scan(returns...)
		if errors.Is(err, sql.ErrNoRows) {
			err = nil // Postgres doesn't return anything when there's no update
		}
	} else {
		_, err = exec.ExecContext(ctx, cache.query, vals...)
	}
	if err != nil {
		return errors.Wrap(err, "models: unable to upsert user_position")
	}

	if !cached {
		userPositionUpsertCacheMut.Lock()
		userPositionUpsertCache[key] = cache
		userPositionUpsertCacheMut.Unlock()
	}

	return o.doAfterUpsertHooks(ctx, exec)
}

// Delete deletes a single UserPosition record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *UserPosition) Delete(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if o == nil {
		return 0, errors.New("models: no UserPosition provided for delete")
	}

	if err := o.doBeforeDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), userPositionPrimaryKeyMapping)
	sql := "DELETE FROM \"user_position\" WHERE \"pubkey\"=$1"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete from user_position")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by delete for user_position")
	}

	if err := o.doAfterDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	return rowsAff, nil
}

// DeleteAll deletes all matching rows.
func (q userPositionQuery) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if q.Query == nil {
		return 0, errors.New("models: no userPositionQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from user_position")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for user_position")
	}

	return rowsAff, nil
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o UserPositionSlice) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if len(o) == 0 {
		return 0, nil
	}

	if len(userPositionBeforeDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doBeforeDeleteHooks(ctx, exec); err != nil {
				return 0, err
			}
		}
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), userPositionPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM \"user_position\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, userPositionPrimaryKeyColumns, len(o))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from userPosition slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for user_position")
	}

	if len(userPositionAfterDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterDeleteHooks(ctx, exec); err != nil {
				return 0, err
			}
		}
	}

	return rowsAff, nil
}

// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *UserPosition) Reload(ctx context.Context, exec boil.ContextExecutor) error {
	ret, err := FindUserPosition(ctx, exec, o.Pubkey)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *UserPositionSlice) ReloadAll(ctx context.Context, exec boil.ContextExecutor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := UserPositionSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), userPositionPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT \"user_position\".* FROM \"user_position\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, userPositionPrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(ctx, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "models: unable to reload all in UserPositionSlice")
	}

	*o = slice

	return nil
}

// UserPositionExists checks if the UserPosition row exists.
func UserPositionExists(ctx context.Context, exec boil.ContextExecutor, pubkey string) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from \"user_position\" where \"pubkey\"=$1 limit 1)"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, pubkey)
	}
	row := exec.QueryRowContext(ctx, sql, pubkey)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "models: unable to check if user_position exists")
	}

	return exists, nil
}
