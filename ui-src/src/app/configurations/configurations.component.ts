import { Component, OnInit } from "@angular/core";
import { HttpModule } from "@angular/http";
import { Http } from "@angular/http";
import { PostService } from "../post.service";

@Component({
  selector: "app-configurations",
  templateUrl: "./configurations.component.html",
  styleUrls: ["./configurations.component.css"]
})
export class ConfigurationsComponent {
  exampleSchema = {
    type: "object",
    properties: {
      applicationName: { type: "string" },
      targetHost: { type: "string" },
      targetPort: {
        type: "string",
        minimum: 1,
        maximum: 65535
      },
      memoryEndpoint: { type: "string" },
      numIterations: { type: "integer", minimum: 0 },
      concurrentUsers: { type: "integer", minimum: 0 },
      memoryVariance: {
        type: "number",
        minimum: 0,
        maximum: 100
      },
      serviceVariance: {
        type: "number",
        minimum: 0,
        maximum: 100
      },
      testSuite: {
        type: "string",
        enum: ["Default-1", "Default-2", "Default-3"]
      },
      requestDelay: { type: "integer", minimum: 0 },
      tpsFrequency: { type: "integer", minimum: 0 },
      rampUsers: { type: "integer", minimum: 0 },
      rampDelay: { type: "integer", minimum: 0 },
      testCaseDirectory: { type: "string" },
      testSuiteDirectory: { type: "string" },
      baseStatsOutputDirectory: { type: "string" },
      reportOutputDirectory: { type: "string" }
    },
    required: [
      "applicationName",
      "targetHost",
      "targetPort",
      "numIterations",
      "concurrentUsers",
      "memoryVariance",
      "serviceVariance",
      "testSuite",
      "requestDelay",
      "tpsFrequency",
      "rampUsers",
      "rampDelay",
      "testCaseDirectory",
      "testSuiteDirectory",
      "baseStatsOutputDirectory",
      "reportOutputDirectory"
    ]
  };

  exampleData = {
    allowablePeakMemoryVariance: 15,
    allowableServiceResponseTimeVariance: 15
  };

  layout: [
    {
      type: "flex";
      "flex-flow": "row wrap";
      items: ["applicationName", "targetPort"];
    }
  ];
  displayData: any = null;
  constructor(private postService: PostService) {}

  exampleOnSubmitFn(form) {
    // this.displayData = form;
    this.postService.addConfig(form);
  }
}
