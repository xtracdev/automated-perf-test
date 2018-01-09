import { Component, OnInit } from "@angular/core";
import { FormControl, FormGroup, Validators } from "@angular/forms";

@Component({
  selector: "app-test-cases",
  templateUrl: "./test-cases.component.html",
  styleUrls: ["./test-cases.component.css"]
})
export class TestCasesComponent implements OnInit {
  formGroup: FormGroup;

  constructor() {}

  ngOnInit() {}
  onSubmit() {
    console.log(this.formGroup);
  }
}
