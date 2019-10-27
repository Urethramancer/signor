# Stringer
Stringer extends strings.Builder with some convenience functions.

## Usage

### Just strings

```go
	s := stringer.New()
	s.WriteStrings("one", "two")
```

### Mixed data

#### Strings and numbers

```go
	s := stringer.New()
	s.WriteI("one", 2, 3.01)
```

#### Slices

```go
	s := stringer.New()
	s.WriteI([]int{1, 2, 3})
```

#### Slices with comma-separation

```go
	s := stringer.New().SetSliceComma(true)
	s.WriteI([]int{1, 2, 3})
```

#### Maps

```go
	s := stringer.New()
	s.WriteI(map[int]string{1: "one", 2: "two", 3: "three"})
```

#### Maps with comma-separation

```go
	s := stringer.New().SetMapComma(true)
	s.WriteI(map[int]string{1: "one", 2: "two", 3: "three"})
```

#### Maps with alternative key-value separator

```go
	s := stringer.New().SetMapComma(true).SetEquals(':')
	s.WriteI(map[int]string{1: "one", 2: "two", 3: "three"})
```

#### Booleans

```go
	s := stringer.New()
	s.WriteI(true)
```
