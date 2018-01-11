import { Component, OnInit, Output, EventEmitter, Input } from "@angular/core";
import { TestCriteria } from "../test-criteria";

@Component({
  selector: "app-test-criteria",
  templateUrl: "./test-criteria.component.html",
  styleUrls: ["./test-criteria.component.css"]
})
export class TestCriteriaComponent {
  @Input() testCriteria: TestCriteria;

  @Output() testCPropertiesChange = new EventEmitter<TestCriteria>();

  formChanged() {
    this.testCPropertiesChange.emit(this.testCriteria);
  }
}
