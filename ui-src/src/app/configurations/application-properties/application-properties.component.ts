import { Component, OnInit, Output, EventEmitter, Input } from "@angular/core";
import { ApplicationProperties } from "../application-properties";
import { FormControl, FormGroup, Validator } from "@angular/forms";
@Component({
  selector: "app-application-properties",
  templateUrl: "./application-properties.component.html",
  styleUrls: ["./application-properties.component.css"]
})
export class ApplicationPropertiesComponent {
  @Input() appProperties: ApplicationProperties;

  @Output() appPropertiesChange = new EventEmitter<ApplicationProperties>();

  formChanged() {
    this.appPropertiesChange.emit(this.appProperties);
  }
}
