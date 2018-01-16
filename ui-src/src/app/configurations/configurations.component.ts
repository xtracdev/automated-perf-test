import { Component, OnInit } from "@angular/core";

@Component({
  selector: "app-configurations",
  templateUrl: "./configurations.component.html",
  styleUrls: ["./configurations.component.css"]
})
export class ConfigurationsComponent {
  exampleSchema = {
    "type": "object",
    "required": ["applicationName", "targetHost","configFilePath" ,"targetPort", "memoryVariance",
      "numIterations", "serviceVariance", "testCaseDirectory",
      "testSuiteDirectory", "baseStatsDirectory", "reportDirectory", "concurrentUsers",
      "testSuite", "requestDelay", "tpsFrequency", "rampUsers", "rampDelay"],
    "properties":
      {
        "configFilePath": {"type": "string" },
        "applicationName": { "type": "string" },
        "targetHost": { "type": "string" },
        "targetPort": { "type": "string" },
        "numIterations": { "type": "integer", "minimum": 0 },
        "memoryVariance": { "type": "number", "minimum": 0 },
        "serviceVariance": { "type": "number", "minimum": 0 },
        "testCaseDirectory": { "type": "string" },
        "testSuiteDirectory": { "type": "string" },
        "baseStatsDirectory": { "type": "string" },
        "reportDirectory": { "type": "string" },
        "concurrentUsers": { "type": "integer", "minimum": 0 },
        "testSuite": { "type": "string" },
        "memoryEndpoint": { "type": "string" },
        "requestDelay": { "type": "integer", "minimum": 0 },
        "tpsFrequency": { "type": "integer", "minimum": 0 },
        "rampUsers": { "type": "integer", "minimum": 0 },
        "rampDelay": { "type": "integer", "minimum": 0 }
      }
  };

  displayData: any = null;

  exampleOnSubmitFn(formData) {
    this.displayData = formData;
    alert('hi Colm');
  }
}
