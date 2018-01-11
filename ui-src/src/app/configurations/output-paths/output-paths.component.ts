import { Component, OnInit, Input, Output, EventEmitter } from "@angular/core";
import { OutputPaths } from "../output-paths";
@Component({
  selector: "app-output-paths",
  templateUrl: "./output-paths.component.html",
  styleUrls: ["./output-paths.component.css"]
})
export class OutputPathsComponent {
  @Input() outputPaths: OutputPaths;

  @Output() outputsPropertiesChange = new EventEmitter<OutputPaths>();

  formChanged() {
    this.outputsPropertiesChange.emit(this.outputPaths);
  }
}
