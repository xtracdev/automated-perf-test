import { Component, OnInit } from "@angular/core";
import { FormControl, FormGroup, Validators } from "@angular/forms";

@Component({
  selector: "app-configurations",
  templateUrl: "./configurations.component.html",
  styleUrls: ["./configurations.component.css"]
})
export class ConfigurationsComponent implements OnInit {
  formGroup: FormGroup;
  constructor() {}

  ngOnInit() {}
  onSubmit() {
    console.log(this.formGroup);
  }
}
