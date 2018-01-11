import { Component, OnInit } from "@angular/core";
import { TestCriteria } from "../test-criteria";

@Component({
  selector: "app-test-criteria",
  templateUrl: "./test-criteria.component.html",
  styleUrls: ["./test-criteria.component.css"]
})
export class TestCriteriaComponent {
  model = new TestCriteria(null, null, null, null, null, null, null, null, "");

  submitted = false;

  onSubmit() {
    this.submitted = true;
  }
}
