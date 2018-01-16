import { Component, OnInit } from "@angular/core";

@Component({
  selector: "app-configurations",
  templateUrl: "./configurations.component.html",
  styleUrls: ["./configurations.component.css"]
})
export class ConfigurationsComponent {
  exampleSchema = {
    "type": "object",
    "required": [
      "apiName",
      "targetHost",
      "targetPort",
      "numIterations",
      "concurrentUsers",
      "allowablePeakMemoryVariance",
      "allowableServiceResponseTimeVariance",
      "testSuite",
      "requestDelay",
      "TPSFreq",
      "rampUsers",
      "rampDelay",
      "testCaseDir",
      "testSuiteDir",
      "baseStatsOutputDir",
      "reportOutputDir"
    ],
    "properties": {
      "apiName": { "type": "string" },
      "targetHost": { "type": "string" },
      " targetPort": {
        "type": "string",
        "minimum": 1,
        "maximum": 65535
      },
      "memoryEndpoint": { "type": "string" },
      "numIterations": { "type": "integer", "minimum": 0 },
      "concurrentUsers": { "type": "integer", "minimum": 0 },
      "allowablePeakMemoryVariance": {
        " type": "integer",
        "minimum": 0,
        "maximum": 100
      },
      "allowableServiceResponseTimeVariance": {
        "type": "integer",
        "minimum": 0,
        " maximum": 100
      },
      "testSuite": {
        "type": "string",
        "enum": ["Default-1", "Default-2", "Default-3"]
      },
      "requestDelay": { "type": "integer", "minimum": 0 },
      "TPSFreq": { "type": "integer", "minimum": 0 },
      "rampUsers": { "type": "integer", "minimum": 0 },
      "rampDelay": { "type": "integer", "minimum": 0 },
      "testCaseDir": { "type": "string" },
      "testSuiteDir": { "type": "string" },
      "baseStatsOutputDir": { "type": "string" },
      "reportOutputDir": { "type": "string" }
    }
  };
  

  displayData: any = null;

  exampleOnSubmitFn(formData) {
    this.displayData = formData;
    alert('hi Colm');
  }
}
