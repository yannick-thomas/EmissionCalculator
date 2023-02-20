from xhtml2pdf import pisa
import subprocess
from jinja2 import Environment, FileSystemLoader

class PrintHandler():
    
    def __init__(self, results : list) -> None:
        self.results = results

    def build_template(self): 
        environment = Environment(loader=FileSystemLoader("templates/"))
        template = environment.get_template("print_template.html")

        return template.render(
            emission_component_result = str(self.results[2]).replace(".", ","),
            emissions = self.results[1],
            energy_content = self.results[3]
        )
    
    def convert_template(self):
        with open("label.pdf", "w+b") as output_file:
            pisa.CreatePDF(self.build_template(), dest=output_file, pagesize=(210, 203))

    def print_result(self):
        self.convert_template()
        subprocess.call(['rundll32','shimgvw.dll,ImageView_PrintTo', 'label.pdf'])