import { Component, OnInit } from "@angular/core";

@Component({
  selector: "app-test-suites",
  templateUrl: "./test-suites.component.html",
  styleUrls: ["./test-suites.component.css"]
})
export class TestSuitesComponent {

  testSchema = {
    type: "object",
    properties: {
      name: { type: "string" },
      testStrategy: { type: "string" },
      description: { type: "string" }

    },
    layout: [
      {

        submit: "hidden"
      }
    ],
  }

  onAdd() { }
  onDelete() { }
  onCancel() { }
  onSave() { }

}
