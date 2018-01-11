import { Component, OnInit } from "@angular/core";
import { FormControl, FormGroup, Validators } from "@angular/forms";
import { ApplicationProperties } from "./application-properties";
import { TestCriteria } from "./test-criteria";
import { OutputPaths } from "./output-paths";
@Component({
  selector: "app-configurations",
  templateUrl: "./configurations.component.html",
  styleUrls: ["./configurations.component.css"]
})
export class ConfigurationsComponent {
  applicationProperties: ApplicationProperties = new ApplicationProperties();
  testCrit: TestCriteria = new TestCriteria();
  outPaths: OutputPaths = new OutputPaths();
  submit() {
    alert(JSON.stringify(this.applicationProperties));
    alert(JSON.stringify(this.testCrit));
    alert(JSON.stringify(this.outPaths));
  }
}
