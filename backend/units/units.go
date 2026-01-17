package units

import (
	"fmt"
	"math"
	"sort"
	"strings"

	"github.com/cobyabrahams/hungr/models"
)

const IntegerTolerance = 0.01

type UnitCategory string

const (
	CategoryVolume UnitCategory = "volume"
	CategoryMass   UnitCategory = "mass"
	CategoryCount  UnitCategory = "count"
)

type DerivedUnit struct {
	ToBase       float64
	Name         string
	Abbrev       string
	PluralName   string
	PluralAbbrev string
	NoDisplay    bool // If true, unit can be parsed but won't be chosen for display
}

// VolumeUnits (base: ml)
var VolumeUnits = map[string]DerivedUnit{
	"ml": {ToBase: 1, Name: "milliliter", Abbrev: "ml", PluralName: "milliliters"},
	"l":  {ToBase: 1000, Name: "liter", Abbrev: "l", PluralName: "liters"},

	"tsp":      {ToBase: 4.92892, Name: "teaspoon", Abbrev: "tsp", PluralName: "teaspoons"},
	"half_tsp": {ToBase: 2.46446, Name: "half teaspoon", Abbrev: "½ tsp", PluralName: "half teaspoons"},
	"qtr_tsp":  {ToBase: 1.23223, Name: "quarter teaspoon", Abbrev: "¼ tsp", PluralName: "quarter teaspoons"},
	"eighth_tsp": {ToBase: 0.616115, Name: "eighth teaspoon", Abbrev: "⅛ tsp", PluralName: "eighth teaspoons"},
	"tbsp":   {ToBase: 14.7868, Name: "tablespoon", Abbrev: "tbsp", PluralName: "tablespoons"},
	"fl_oz":  {ToBase: 29.5735, Name: "fluid ounce", Abbrev: "fl oz", PluralName: "fluid ounces", NoDisplay: true},
	"jigger": {ToBase: 44.3603, Name: "jigger", Abbrev: "jigger", PluralName: "jiggers", NoDisplay: true},
	"cup":      {ToBase: 236.588, Name: "cup", Abbrev: "cup", PluralName: "cups"},
	"half_cup": {ToBase: 118.294, Name: "half cup", Abbrev: "½ cup", PluralName: "half cups"},
	"gill":     {ToBase: 118.294, Name: "gill", Abbrev: "gill", PluralName: "gills", NoDisplay: true},
	"qtr_cup":  {ToBase: 59.147, Name: "quarter cup", Abbrev: "¼ cup", PluralName: "quarter cups"},
	//	"pt":     {ToBase: 473.176, Name: "pint", Abbrev: "pt", PluralName: "pints"},
	"qt":  {ToBase: 946.353, Name: "quart", Abbrev: "qt", PluralName: "quarts"},
	"gal": {ToBase: 3785.41, Name: "gallon", Abbrev: "gal", PluralName: "gallons"},

	"drop":    {ToBase: 0.05, Name: "drop", Abbrev: "drop", PluralName: "drops", NoDisplay: true},
	"smidgen": {ToBase: 0.115522, Name: "smidgen", Abbrev: "smidgen", PluralName: "smidgens", NoDisplay: true},
	"pinch":   {ToBase: 0.231043, Name: "pinch", Abbrev: "pinch", PluralName: "pinches"},
	"dash":    {ToBase: 0.462086, Name: "dash", Abbrev: "dash", PluralName: "dashes"},

	"imp_tsp":    {ToBase: 5.91939, Name: "imperial teaspoon", Abbrev: "imp tsp", PluralName: "imperial teaspoons"},
	"imp_tbsp":   {ToBase: 17.7582, Name: "imperial tablespoon", Abbrev: "imp tbsp", PluralName: "imperial tablespoons"},
	"imp_fl_oz":  {ToBase: 28.4131, Name: "imperial fluid ounce", Abbrev: "imp fl oz", PluralName: "imperial fluid ounces"},
	"imp_cup":    {ToBase: 284.131, Name: "imperial cup", Abbrev: "imp cup", PluralName: "imperial cups"},
	"imp_pt":     {ToBase: 568.261, Name: "imperial pint", Abbrev: "imp pt", PluralName: "imperial pints"},
	"imp_qt":     {ToBase: 1136.52, Name: "imperial quart", Abbrev: "imp qt", PluralName: "imperial quarts"},
	"imp_gal":    {ToBase: 4546.09, Name: "imperial gallon", Abbrev: "imp gal", PluralName: "imperial gallons"},
}

// MassUnits (base: mg)
var MassUnits = map[string]DerivedUnit{
	"mcg": {ToBase: 0.001, Name: "microgram", Abbrev: "mcg", PluralName: "micrograms"},
	"mg":  {ToBase: 1, Name: "milligram", Abbrev: "mg", PluralName: "milligrams"},
	"g":   {ToBase: 1000, Name: "gram", Abbrev: "g", PluralName: "grams"},
	"kg":  {ToBase: 1000000, Name: "kilogram", Abbrev: "kg", PluralName: "kilograms"},

	"gr":    {ToBase: 64.79891, Name: "grain", Abbrev: "gr", PluralName: "grains"},
	"dr":    {ToBase: 1771.8452, Name: "dram", Abbrev: "dr", PluralName: "drams"},
	"oz":    {ToBase: 28349.5, Name: "ounce", Abbrev: "oz", PluralName: "ounces"},
	"lb":    {ToBase: 453592, Name: "pound", Abbrev: "lb", PluralName: "pounds"},
	"stone": {ToBase: 6350293, Name: "stone", Abbrev: "st", PluralName: "stone"},

	"oz_t": {ToBase: 31103.5, Name: "troy ounce", Abbrev: "oz t", PluralName: "troy ounces"},
	"dwt":  {ToBase: 1555.17, Name: "pennyweight", Abbrev: "dwt", PluralName: "pennyweights"},
}

type Quantity struct {
	Value    float64
	Unit     string
	Category UnitCategory
}

var sortedVolumeUnits []string
var sortedMassUnits []string

func init() {
	sortedVolumeUnits = make([]string, 0, len(VolumeUnits))
	for k := range VolumeUnits {
		sortedVolumeUnits = append(sortedVolumeUnits, k)
	}
	sort.Slice(sortedVolumeUnits, func(i, j int) bool {
		return VolumeUnits[sortedVolumeUnits[i]].ToBase > VolumeUnits[sortedVolumeUnits[j]].ToBase
	})

	sortedMassUnits = make([]string, 0, len(MassUnits))
	for k := range MassUnits {
		sortedMassUnits = append(sortedMassUnits, k)
	}
	sort.Slice(sortedMassUnits, func(i, j int) bool {
		return MassUnits[sortedMassUnits[i]].ToBase > MassUnits[sortedMassUnits[j]].ToBase
	})
}

func GetCategoryForIngredientUnit(u models.IngredientUnit) UnitCategory {
	switch u {
	case models.UnitML:
		return CategoryVolume
	case models.UnitMG:
		return CategoryMass
	case models.UnitCount:
		return CategoryCount
	default:
		return ""
	}
}

func ToBaseUnit(value float64, unit string) (float64, UnitCategory, error) {
	if u, ok := VolumeUnits[unit]; ok {
		return value * u.ToBase, CategoryVolume, nil
	}
	if u, ok := MassUnits[unit]; ok {
		return value * u.ToBase, CategoryMass, nil
	}
	if unit == "count" {
		return value, CategoryCount, nil
	}
	return 0, "", fmt.Errorf("unknown unit: %s", unit)
}

func FromBaseUnit(baseValue float64, category UnitCategory, targetUnit string) (float64, error) {
	switch category {
	case CategoryVolume:
		if u, ok := VolumeUnits[targetUnit]; ok {
			return baseValue / u.ToBase, nil
		}
		return 0, fmt.Errorf("unknown volume unit: %s", targetUnit)
	case CategoryMass:
		if u, ok := MassUnits[targetUnit]; ok {
			return baseValue / u.ToBase, nil
		}
		return 0, fmt.Errorf("unknown mass unit: %s", targetUnit)
	case CategoryCount:
		return baseValue, nil
	default:
		return 0, fmt.Errorf("unknown category: %s", category)
	}
}

func Convert(value float64, from, to string) (float64, error) {
	baseValue, category, err := ToBaseUnit(value, from)
	if err != nil {
		return 0, err
	}

	var toCategory UnitCategory
	if _, ok := VolumeUnits[to]; ok {
		toCategory = CategoryVolume
	} else if _, ok := MassUnits[to]; ok {
		toCategory = CategoryMass
	} else {
		return 0, fmt.Errorf("unknown target unit: %s", to)
	}

	if category != toCategory {
		return 0, fmt.Errorf("cannot convert between %s (%s) and %s (%s)", from, category, to, toCategory)
	}

	return FromBaseUnit(baseValue, category, to)
}

func isNearInteger(value float64, tolerance float64) bool {
	rounded := math.Round(value)
	if rounded == 0 {
		return false
	}
	diff := math.Abs(value - rounded)
	return diff/rounded <= tolerance
}

// FindBestIntegerUnit finds the largest unit where the value converts to within
// IntegerTolerance of an integer. Returns base unit if no match found.
func FindBestIntegerUnit(baseValue float64, category UnitCategory) Quantity {
	return FindBestIntegerUnitWithTolerance(baseValue, category, IntegerTolerance)
}

func FindBestIntegerUnitWithTolerance(baseValue float64, category UnitCategory, tolerance float64) Quantity {
	if category == CategoryCount {
		return Quantity{Value: baseValue, Unit: "count", Category: CategoryCount}
	}

	var sortedUnits []string
	var unitMap map[string]DerivedUnit
	var baseUnit string

	switch category {
	case CategoryVolume:
		sortedUnits = sortedVolumeUnits
		unitMap = VolumeUnits
		baseUnit = "ml"
	case CategoryMass:
		sortedUnits = sortedMassUnits
		unitMap = MassUnits
		baseUnit = "mg"
	default:
		return Quantity{Value: baseValue, Unit: "", Category: category}
	}

	for _, unitKey := range sortedUnits {
		unit := unitMap[unitKey]
		if unit.NoDisplay {
			continue
		}
		converted := baseValue / unit.ToBase
		if converted >= 1 && isNearInteger(converted, tolerance) {
			return Quantity{
				Value:    math.Round(converted),
				Unit:     unitKey,
				Category: category,
			}
		}
	}

	return Quantity{Value: baseValue, Unit: baseUnit, Category: category}
}

func Format(q Quantity) string {
	if q.Category == CategoryCount {
		if q.Value == 1 {
			return "1"
		}
		return fmt.Sprintf("%.0f", q.Value)
	}

	var unit DerivedUnit
	var ok bool
	switch q.Category {
	case CategoryVolume:
		unit, ok = VolumeUnits[q.Unit]
	case CategoryMass:
		unit, ok = MassUnits[q.Unit]
	}

	if !ok {
		return fmt.Sprintf("%.2f %s", q.Value, q.Unit)
	}

	var formatted string
	if q.Value == float64(int(q.Value)) {
		formatted = fmt.Sprintf("%d", int(q.Value))
	} else if q.Value < 0.1 {
		formatted = fmt.Sprintf("%.3f", q.Value)
	} else if q.Value < 10 {
		formatted = fmt.Sprintf("%.2f", q.Value)
	} else {
		formatted = fmt.Sprintf("%.1f", q.Value)
	}

	abbrev := unit.Abbrev
	if q.Value != 1 && unit.PluralAbbrev != "" {
		abbrev = unit.PluralAbbrev
	}

	return fmt.Sprintf("%s %s", formatted, abbrev)
}

func FormatBest(baseValue float64, category UnitCategory) string {
	q := FindBestIntegerUnit(baseValue, category)
	return Format(q)
}

func ParseUnit(s string) (string, UnitCategory, error) {
	s = strings.ToLower(strings.TrimSpace(s))

	if _, ok := VolumeUnits[s]; ok {
		return s, CategoryVolume, nil
	}

	if _, ok := MassUnits[s]; ok {
		return s, CategoryMass, nil
	}

	volumeAliases := map[string]string{
		"milliliter":            "ml",
		"milliliters":           "ml",
		"millilitre":            "ml",
		"millilitres":           "ml",
		"centiliter":            "cl",
		"centiliters":           "cl",
		"centilitre":            "cl",
		"centilitres":           "cl",
		"deciliter":             "dl",
		"deciliters":            "dl",
		"decilitre":             "dl",
		"decilitres":            "dl",
		"liter":                 "l",
		"liters":                "l",
		"litre":                 "l",
		"litres":                "l",
		"teaspoon":              "tsp",
		"teaspoons":             "tsp",
		"half teaspoon":         "half_tsp",
		"half teaspoons":        "half_tsp",
		"1/2 tsp":               "half_tsp",
		"1/2 teaspoon":          "half_tsp",
		"1/2 teaspoons":         "half_tsp",
		"½ tsp":                 "half_tsp",
		"½ teaspoon":            "half_tsp",
		"quarter teaspoon":      "qtr_tsp",
		"quarter teaspoons":     "qtr_tsp",
		"1/4 tsp":               "qtr_tsp",
		"1/4 teaspoon":          "qtr_tsp",
		"1/4 teaspoons":         "qtr_tsp",
		"¼ tsp":                 "qtr_tsp",
		"¼ teaspoon":            "qtr_tsp",
		"eighth teaspoon":       "eighth_tsp",
		"eighth teaspoons":      "eighth_tsp",
		"1/8 tsp":               "eighth_tsp",
		"1/8 teaspoon":          "eighth_tsp",
		"1/8 teaspoons":         "eighth_tsp",
		"⅛ tsp":                 "eighth_tsp",
		"⅛ teaspoon":            "eighth_tsp",
		"tablespoon":            "tbsp",
		"tablespoons":           "tbsp",
		"tbsps":                 "tbsp",
		"fluid ounce":           "fl_oz",
		"fluid ounces":          "fl_oz",
		"fl oz":                 "fl_oz",
		"floz":                  "fl_oz",
		"cups":                  "cup",
		"half cup":              "half_cup",
		"half cups":             "half_cup",
		"1/2 cup":               "half_cup",
		"1/2 cups":              "half_cup",
		"½ cup":                 "half_cup",
		"½ cups":                "half_cup",
		"quarter cup":           "qtr_cup",
		"quarter cups":          "qtr_cup",
		"1/4 cup":               "qtr_cup",
		"1/4 cups":              "qtr_cup",
		"¼ cup":                 "qtr_cup",
		"¼ cups":                "qtr_cup",
		"pint":                  "pt",
		"pints":                 "pt",
		"quart":                 "qt",
		"quarts":                "qt",
		"gallon":                "gal",
		"gallons":               "gal",
		"drops":                 "drop",
		"dashes":                "dash",
		"pinches":               "pinch",
		"smidgens":              "smidgen",
		"jiggers":               "jigger",
		"shot":                  "jigger",
		"shots":                 "jigger",
		"gills":                 "gill",
		"australian tablespoon": "au_tbsp",
	}

	massAliases := map[string]string{
		"microgram":    "mcg",
		"micrograms":   "mcg",
		"µg":           "mcg",
		"ug":           "mcg",
		"milligram":    "mg",
		"milligrams":   "mg",
		"centigram":    "cg",
		"centigrams":   "cg",
		"decigram":     "dg",
		"decigrams":    "dg",
		"gram":         "g",
		"grams":        "g",
		"decagram":     "dag",
		"decagrams":    "dag",
		"hectogram":    "hg",
		"hectograms":   "hg",
		"kilogram":     "kg",
		"kilograms":    "kg",
		"kilo":         "kg",
		"kilos":        "kg",
		"grain":        "gr",
		"grains":       "gr",
		"dram":         "dr",
		"drams":        "dr",
		"ounce":        "oz",
		"ounces":       "oz",
		"pound":        "lb",
		"pounds":       "lb",
		"lbs":          "lb",
		"stones":       "stone",
		"st":           "stone",
		"troy ounce":   "oz_t",
		"troy ounces":  "oz_t",
		"pennyweight":  "dwt",
		"pennyweights": "dwt",
	}

	if unit, ok := volumeAliases[s]; ok {
		return unit, CategoryVolume, nil
	}

	if unit, ok := massAliases[s]; ok {
		return unit, CategoryMass, nil
	}

	countAliases := []string{"count", "piece", "pieces", "item", "items", "each", "ea", "unit", "units"}
	for _, alias := range countAliases {
		if s == alias {
			return "count", CategoryCount, nil
		}
	}

	return "", "", fmt.Errorf("unknown unit: %s", s)
}

func ListUnitsForCategory(cat UnitCategory) []string {
	switch cat {
	case CategoryVolume:
		units := make([]string, 0, len(VolumeUnits))
		for k := range VolumeUnits {
			units = append(units, k)
		}
		return units
	case CategoryMass:
		units := make([]string, 0, len(MassUnits))
		for k := range MassUnits {
			units = append(units, k)
		}
		return units
	case CategoryCount:
		return []string{"count"}
	default:
		return nil
	}
}

func ScaleQuantity(q Quantity, factor float64) Quantity {
	return Quantity{
		Value:    q.Value * factor,
		Unit:     q.Unit,
		Category: q.Category,
	}
}

func SumQuantities(quantities []Quantity) (float64, UnitCategory, error) {
	if len(quantities) == 0 {
		return 0, "", fmt.Errorf("no quantities to sum")
	}

	category := quantities[0].Category
	var total float64

	for _, q := range quantities {
		if q.Category != category {
			return 0, "", fmt.Errorf("cannot sum different categories: %s and %s", category, q.Category)
		}

		if category == CategoryCount {
			total += q.Value
			continue
		}

		baseValue, _, err := ToBaseUnit(q.Value, q.Unit)
		if err != nil {
			return 0, "", err
		}
		total += baseValue
	}

	return total, category, nil
}

func GetDerivedUnit(unitKey string) (DerivedUnit, UnitCategory, error) {
	if u, ok := VolumeUnits[unitKey]; ok {
		return u, CategoryVolume, nil
	}
	if u, ok := MassUnits[unitKey]; ok {
		return u, CategoryMass, nil
	}
	return DerivedUnit{}, "", fmt.Errorf("unknown unit: %s", unitKey)
}

type ParsedIngredient struct {
	Quantity       float64
	Unit           string
	Category       UnitCategory
	IngredientName string
}

// ParseIngredientString parses strings like "2 cups flour" or "1/2 tsp salt"
// Also handles ingredients without quantities like "salt to taste" or "avocado oil"
// Returns quantity, unit key, category, and ingredient name
func ParseIngredientString(s string) (ParsedIngredient, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return ParsedIngredient{}, fmt.Errorf("empty ingredient string")
	}

	parts := strings.Fields(s)
	if len(parts) < 1 {
		return ParsedIngredient{}, fmt.Errorf("ingredient string too short: %q", s)
	}

	// Try to parse quantity (first part) - handle fractions like "1/2"
	// If the first word doesn't look like a number (no digits), treat whole string as ingredient
	firstWord := parts[0]
	hasDigit := false
	for _, c := range firstWord {
		if c >= '0' && c <= '9' {
			hasDigit = true
			break
		}
	}

	if !hasDigit {
		// First word isn't a number - treat entire string as ingredient name
		// with quantity 1 and category count
		return ParsedIngredient{
			Quantity:       1,
			Unit:           "count",
			Category:       CategoryCount,
			IngredientName: s,
		}, nil
	}

	quantity, err := parseQuantity(firstWord)
	if err != nil {
		return ParsedIngredient{}, fmt.Errorf("invalid quantity %q: %w", firstWord, err)
	}

	// If only one part and it's a number, that's not valid
	if len(parts) < 2 {
		return ParsedIngredient{}, fmt.Errorf("ingredient string too short: %q", s)
	}

	// Try to find unit starting from second part
	// Handle multi-word units like "fl oz"
	unitEndIdx := 1
	var unitKey string
	var category UnitCategory

	// Try two-word unit first (e.g., "fl oz")
	if len(parts) >= 3 {
		twoWord := parts[1] + " " + parts[2]
		if key, cat, err := ParseUnit(twoWord); err == nil {
			unitKey = key
			category = cat
			unitEndIdx = 3
		}
	}

	// Try single-word unit
	if unitKey == "" {
		if key, cat, err := ParseUnit(parts[1]); err == nil {
			unitKey = key
			category = cat
			unitEndIdx = 2
		}
	}

	// If no unit found, assume count
	if unitKey == "" {
		unitKey = "count"
		category = CategoryCount
		unitEndIdx = 1
	}

	// Rest is ingredient name
	ingredientName := strings.Join(parts[unitEndIdx:], " ")
	if ingredientName == "" {
		return ParsedIngredient{}, fmt.Errorf("no ingredient name found in %q", s)
	}

	return ParsedIngredient{
		Quantity:       quantity,
		Unit:           unitKey,
		Category:       category,
		IngredientName: ingredientName,
	}, nil
}

func parseQuantity(s string) (float64, error) {
	// Handle fractions like "1/2"
	if strings.Contains(s, "/") {
		parts := strings.Split(s, "/")
		if len(parts) != 2 {
			return 0, fmt.Errorf("invalid fraction")
		}
		var num, denom float64
		if _, err := fmt.Sscanf(parts[0], "%f", &num); err != nil {
			return 0, err
		}
		if _, err := fmt.Sscanf(parts[1], "%f", &denom); err != nil {
			return 0, err
		}
		if denom == 0 {
			return 0, fmt.Errorf("division by zero")
		}
		return num / denom, nil
	}

	var val float64
	if _, err := fmt.Sscanf(s, "%f", &val); err != nil {
		return 0, err
	}
	return val, nil
}
