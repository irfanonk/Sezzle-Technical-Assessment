# Test Results

Date: September 22, 2026

## Backend

Command:

```sh
go test -count=1 ./...
```

Result: **Passed**

```text
?    go-calculator/cmd/server             [no test files]
ok   go-calculator/internal/calculator    0.507s
ok   go-calculator/internal/httpapi       0.796s
```

The calculator domain tests and HTTP API tests passed. The server entry-point
package does not contain tests.


% Coverage report

```text
        go-calculator/cmd/server                coverage: 0.0% of statements
ok      go-calculator/internal/calculator       (cached)        coverage: 100.0% of statements
ok      go-calculator/internal/httpapi  0.471s  coverage: 94.4% of statements
```

## Frontend

Command:

```sh
pnpm test
```

Runtime: Node.js 22

Result: **Passed**

```text
Test Files  3 passed (3)
Tests       19 passed (19)
Duration    2.11s
```

The API client, operand validation, and calculator component test suites
passed.

 % Coverage report from v8
--------------------------|---------|----------|---------|---------|-------------------
File                      | % Stmts | % Branch | % Funcs | % Lines | Uncovered Line #s 
--------------------------|---------|----------|---------|---------|-------------------
All files                 |      96 |    83.56 |   95.83 |   95.83 |                   
 api                      |   87.87 |    52.38 |      90 |    87.5 |                   
  calculator.ts           |     100 |      100 |     100 |     100 |                   
  client.ts               |   93.33 |    81.81 |     100 |   93.33 | 48                
  errors.ts               |   78.57 |       20 |      75 |   78.57 | 41-50             
 components               |     100 |      100 |     100 |     100 |                   
  Calculator.tsx          |     100 |      100 |     100 |     100 |                   
  ErrorNotice.tsx         |     100 |      100 |     100 |     100 |                   
  OperandFields.tsx       |     100 |      100 |     100 |     100 |                   
  OperationSelect.tsx     |     100 |      100 |     100 |     100 |                   
  ResultPanel.tsx         |     100 |      100 |     100 |     100 |                   
 lib                      |     100 |      100 |     100 |     100 |                   
  formatResult.ts         |     100 |      100 |     100 |     100 |                   
 queries                  |    90.9 |        0 |   85.71 |    90.9 |                   
  queryClient.ts          |   83.33 |        0 |   66.66 |   83.33 | 25                
  useCalculateMutation.ts |     100 |      100 |     100 |     100 |                   
  useOperationsQuery.ts   |     100 |      100 |     100 |     100 |                   
 test                     |     100 |      100 |     100 |     100 |                   
  utils.tsx               |     100 |      100 |     100 |     100 |                   
 validation               |     100 |      100 |     100 |     100 |                   
  operands.ts             |     100 |      100 |     100 |     100 |                   
--------------------------|---------|----------|---------|---------|-------------------

