import { Component, OnInit } from "@angular/core";
import { ApplicationProperties } from "../application-properties";
@Component({
  selector: "app-application-properties",
  templateUrl: "./application-properties.component.html",
  styleUrls: ["./application-properties.component.css"]
})
export class ApplicationPropertiesComponent {
  model = new ApplicationProperties("", "", null, "");

  submitted = false;

  onSubmit() {
    this.submitted = true;
  }
}
