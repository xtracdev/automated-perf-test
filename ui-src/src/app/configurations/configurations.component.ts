import { Component, OnInit } from "@angular/core";
import { Http } from "@angular/http";
import { AutomatedUIServices } from "../automated-ui-services";
import { JsonSchemaFormModule } from "angular2-json-schema-form";
import { ToastsManager } from "ng2-toastr/ng2-toastr";
import { ToastOptions } from "ng2-toastr/src/toast-options";
import "rxjs/add/operator/map";

@Component({
  selector: "app-configurations",
  templateUrl: "./configurations.component.html",
  styleUrls: ["./configurations.component.css"]
})
export class ConfigurationsComponent implements OnInit {
  formData = {};
  configPath = undefined;
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
        this.toastr.success("Your data has been save!", "Success!");
      },
      error => { switch (error.status) {
        case 500: {
          this.toastr.error ("An error has occurred. Check the logs.");
                          break;
                      }
                      case 400: {
                       this.toastr.error("Some of the fields do not conform to the schema", "An error occurred");
                          break;
                      }
                      default: {
                        this.toastr.error("Your data did not save.", "An error occurred");
                                    }
                    }
      }
    );
  }

  onCancel() {
    this.configPath = undefined;
    this.formData = undefined;
  }
}
