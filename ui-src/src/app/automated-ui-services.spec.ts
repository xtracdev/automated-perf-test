import { TestBed, inject } from "@angular/core/testing";

import { AutomatedUIServices } from "./automated-ui-services";

describe("AutomatedUIServices", () => {
  beforeEach(() => {
    TestBed.configureTestingModule({
      providers: [AutomatedUIServices]
    });
  });

  it(
    "should be created",
    inject([AutomatedUIServices], (service: AutomatedUIServices) => {
      expect(service).toBeTruthy();
    })
  );
});
// it("should create Json", createJsonFile => {
//   expect(createJsonFile).toBe("hello");
// });

//   describe("Hello world", () => {
//     it("says hello", () => {
//       expect("hello").toEqual("hello");
//     });
//});

// export class MockData {
//   apiName: "UnitTest";
//   targetHost: "UnitTest";
//   targetPort: "UnitTest";
//   numIterations: 1;
//   concurrentUsers: 1;
//   allowablePeakMemoryVariance: 1;
//   allowableServiceResponseTimeVariance: 1;
//   testSuite: "UnitTest";
//   requestDelay: 1;
//   TPSFreq: 1;
//   rampUsers: 1;
//   rampDelay: 1;
//   testCaseDir: "UnitTest";
//   testSuiteDir: "UnitTest";
//   baseStatsOutputDir: "UnitTest";
//   reportOutputDir: "UnitTest";
// }
