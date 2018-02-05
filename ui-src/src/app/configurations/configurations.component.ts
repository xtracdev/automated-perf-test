import { Component, OnInit } from "@angular/core";
import { Http } from "@angular/http";
import { AutomatedUIServices } from "../automated-ui-services";
import { JsonSchemaFormModule } from "angular2-json-schema-form";
import { ToastsManager } from "ng2-toastr/ng2-toastr";
import { ToastOptions } from "ng2-toastr/src/toast-options";
import "rxjs/add/operator/map";
import { DOCUMENT } from "@angular/platform-browser";

@Component({
  selector: "app-configurations",
  templateUrl: "./configurations.component.html",
  styleUrls: ["./configurations.component.css"]
})
export class ConfigurationsComponent implements OnInit {
  formData = {};
  configPath = undefined;
  xmlFileName = undefined;
  fileName = undefined;
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
    this.automatedUIServices.postConfig$(configData, this.configPath).subscribe(
      data => {
        this.toastr.success("Your Data has Been Saved!", "Success!");
      },

      error => {
        switch (error.status) {
          case 500: {
            this.toastr.error("An Error has Occurred!", "Check the logs!");
            break;
          }
          case 409: {
            this.toastr.error("File Already Exists!", "An Error Occurred!");
            break;
          }
          case 400: {
            this.toastr.error(
              "Some of the Fields do not Conform to the Schema!",
              "An Error Occurred!"
            );
            break;
          }
          default: {
            this.toastr.error("Your Data did not Save!", "An Error Occurred!");
          }
        }
      }
    );
  }

  onCancel() {
    this.automatedUIServices
      .getConfig$(this.configPath, this.xmlFileName)
      .subscribe(
        data => {
          this.formData = data;
          this.toastr.success("Previous Data Reloaded!");
        },
        error => {
          this.configPath = undefined;
          this.xmlFileName = undefined;
          this.formData = undefined;
        }
      );
  }

  fileSelector(event) {
    this.fileName = event.srcElement.files[0].name;
    this.xmlFileName = this.fileName;
    this.onGetFile();
  }
  onGetFile() {
    this.automatedUIServices
      .getConfig$(this.configPath, this.xmlFileName)
      .subscribe(
        data => {
          this.formData = data;
          this.toastr.success("Success!");
        },
        error => {
          switch (error.status) {
            case 404: {
              this.toastr.error("File Not Found!", "An Error Occured!");
              break;
            }
            case 400: {
              this.toastr.error(
                "Check Your Field Inputs",
                "An Error Occurred!"
              );
              break;
            }
            case 500: {
              this.toastr.error("An Error has Occurred!", "Check the logs!");
              break;
            }
            default: {
              this.toastr.error(
                "Your Data was Not Retrieved!",
                "An Error Occurred!"
              );
            }
          }
        }
      );
  }
  onUpdate(configData) {
    this.automatedUIServices
      .putConfig$(this.formData, this.configPath, this.xmlFileName)
      .subscribe(
        data => {
          this.toastr.success("Success!");
        },
        error => {
          switch (error.status) {
            case 404: {
              this.toastr.error("File Not Found", "An Error Occured!");
              break;
            }
            case 409: {
              this.toastr.error(
                "File Must be Specified!",
                "An Error Occurred!"
              );
              break;
            }
            case 400: {
              this.toastr.error("Bad Request!", "An Error Occurred!");
              break;
            }
            case 500: {
              this.toastr.error("Internal Server Error!");
              break;
            }
            default: {
              this.toastr.error("File Was Not Updated!", "An Error Occurred!");
            }
          }
        }
      );
  }
}
