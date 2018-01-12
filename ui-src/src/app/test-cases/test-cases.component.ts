import { Component, OnInit } from "@angular/core";

@Component({
  selector: "app-test-cases",
  templateUrl: "./test-cases.component.html",
  styleUrls: ["./test-cases.component.css"]
})
export class TestCasesComponent {
  title = "the Angular JSON Schema Form<br> Bootstrap 4 Seed App";

  exampleSchema = {
    type: "object",
    required: [
      "apiName",
      "targetHost",
      "targetPort",
      "allowablePeakMemoryVariance",
      "numIterations",
      "allowableServiceResponseTimeVariance",
      "testCaseDir",
      "testSuiteDir",
      "baseStatsOutputDir",
      "reportOutputDir",
      "concurrentUsers",
      "testSuite",
      "requestDelay",
      "TPSFreq",
      "rampUsers",
      "rampDelay"
    ],
    properties: {
      apiName: { type: "string" },
      targetHost: { type: "string" },
      targetPort: { type: "string" },
      numIterations: { type: "integer", minimum: 0 },
      allowablePeakMemoryVariance: { type: "number", minimum: 0 },
      allowableServiceResponseTimeVariance: { type: "number", minimum: 0 },
      testCaseDir: { type: "string" },
      testSuiteDir: { type: "string" },
      baseStatsOutputDir: { type: "string" },
      reportOutputDir: { type: "string" },
      concurrentUsers: { type: "integer", minimum: 0 },
      testSuite: { type: "string" },
      memoryEndpoint: { type: "string" },
      requestDelay: { type: "integer", minimum: 0 },
      TPSFreq: { type: "integer", minimum: 0 },
      rampUsers: { type: "integer", minimum: 0 },
      rampDelay: { type: "integer", minimum: 0 }
    }
  };

  displayData: any = null;

  exampleOnSubmitFn(formData) {
    this.displayData = formData;
  }
}
