import {Component, OnInit, ViewContainerRef} from "@angular/core";
import {ToastsManager} from "ng2-toastr";

@Component({
  selector: "app-root",
  templateUrl: "./app.component.html",
  styleUrls: ["./app.component.css"]
})
export class AppComponent {
  title = "Automated Performance Testing";

  constructor(private toastr: ToastsManager, viewContainerReference: ViewContainerRef) {
    this.toastr.setRootViewContainerRef(viewContainerReference);
  }
}
