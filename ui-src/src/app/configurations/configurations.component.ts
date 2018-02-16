import { Component, OnInit } from "@angular/core";
import { HttpClient } from "@angular/common/http";
import { AutomatedUIServices } from "../automated-ui-services";
import { JsonSchemaFormModule } from "angular2-json-schema-form";
import { ToastsManager } from "ng2-toastr/ng2-toastr";
import "rxjs/add/operator/map";

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
  configSchema = { layout: true };

  constructor(
    private automatedUIServices: AutomatedUIServices,
    private toastr: ToastsManager,

    private http: HttpClient
  ) {}

  ngOnInit() {
    this.automatedUIServices
      .getSchema$("assets/schema.json")
      .subscribe((data: any) => {
        this.configSchema = data;
      });
  }

  onSubmit(configData) {
    this.automatedUIServices.postConfig$(configData, this.configPath).subscribe(
      data => {
        this.toastr.success("Your data has been saved!", "Success!");
      },

      error => {
        switch (error.status) {
          case 500: {
            this.toastr.error("An error has occurred!", "Check the logs!");
            break;
          }
          case 409: {
            this.toastr.error("File already exists!", "An error occurred!");
            break;
          }
          case 400: {
            this.toastr.error(
              "Some of the fields do not conform to the schema!",
              "An error occurred!"
            );
            break;
          }
          default: {
            this.toastr.error("Your data did not save!", "An error occurred!");
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
          this.toastr.success("Previous data reloaded!");
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
    this.xmlFileName = this.xmlFileName.substring(
      0,
      this.xmlFileName.length - 4);


    this.onGetFile();
    (<HTMLInputElement>document.getElementById("file")).value ="";
   
  }

  onClearFile(){
    this.xmlFileName = undefined;
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
              this.toastr.error("File not found!", "An error occured!");
              break;
            }
            case 400: {
              this.toastr.error(
                "Check your field inputs" ,
                "An error occurred!"
              );
              break;
            }
            case 500: {
              this.toastr.error("An error has occurred!", "Check the logs!");
              break;
            }
            default: {
              this.toastr.error(
                "Your data was not retrieved!",
                "An error occurred!"
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
              this.toastr.error("File not found", "An error occured!");
              break;
            }
            case 400: {
              this.toastr.error(
                "File must be specified!",
                "An error occurred!"
              );
              break;
            }
            case 500: {
              this.toastr.error("Internal server error!");
              break;
            }
            default: {
              this.toastr.error("File was not updated!", "An error occurred!");
            }
          }
        }
      );
  }
}
