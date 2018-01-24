import { Component, OnInit } from "@angular/core";
import { Http } from "@angular/http";
import { AutomatedUIServices } from "../automated-ui-services";
import { JsonSchemaFormModule } from "angular2-json-schema-form";
import { ToastsManager } from "ng2-toastr/ng2-toastr";
import { ToastOptions } from "ng2-toastr/src/toast-options";
import { Observable } from "rxjs/Observable";
import "rxjs/add/operator/map";
@Component({
  selector: "app-configurations",
  templateUrl: "./configurations.component.html",
  styleUrls: ["./configurations.component.css"]
})
export class ConfigurationsComponent implements OnInit {
  formData = {};
  // needed for layout to load
  configSchema = { layout: "" };

  constructor(
    private automatedUIServices: AutomatedUIServices,
    private toastr: ToastsManager,
    private http: Http
  ) {}

  ngOnInit() {
    this.http
      .get("assets/schema.json")
      .map((data: any) => data.json())
      .subscribe(
        (data: any) => {
          this.configSchema = data;
        },
        err => console.log(err), // error
        () => console.log("Complete") // complete
      );

    this.formData = {
      allowablePeakMemoryVariance: 15,
      allowableServiceResponseTimeVariance: 15
    };
  }

  onSubmit(configData) {
    this.automatedUIServices.postConfig$(configData).subscribe(
      data => {
        this.toastr.success("Your data has been save!", "Success!");
      },
      error => {
        this.toastr.error("Failed to save data, Check the Command Line!");
      }
    );
  }
  onCancel() {
    this.formData = {};
  }
}
