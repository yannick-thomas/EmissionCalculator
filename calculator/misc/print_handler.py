from xhtml2pdf import pisa
from reportlab.graphics.barcode import code128
from reportlab.graphics.barcode import code93
from reportlab.graphics.barcode import code39
from reportlab.graphics.barcode import usps
from reportlab.graphics.barcode import usps4s
from reportlab.graphics.barcode import ecc200datamatrix
import os
from jinja2 import Environment, BaseLoader
import re
import tempfile

class PrintHandler():

    def __init__(self, results : list) -> None:
        self.tempfile_path = os.path.join(tempfile.gettempdir(), "label.pdf")
        self.results = results

    def build_template(self): 
        template = Environment(loader=BaseLoader).from_string("""
            <!DOCTYPE html>
            <html lang="en">
            <head>
                <style>
                    @page {
                        size: 2481px 2398px landscape;
                        margin-left: 239.8;
                        margin-top: 744.3;
                    }
                    p {
                        font-size: 30px;
                    }
                </style>
            </head>
            <body>
                <p><b>enth. CO2 Abgabe {{ emission_component_result }} EUR, Brutto</b> / CO2 Abgabe pro t 30,00 EUR, Netto</p>
                <p>CO2 Kg der Lieferung {{ emissions }} / CO2 Kg pro kWh {{ co2_per_kwh }} / kWh der Lieferung {{ energy_content }}</p>
            </body>
            </html>
        """)

        return template.render(
            emission_component_result = self.results[2],
            emissions = self.results[1],
            energy_content = self.results[3],
            co2_per_kwh = self.results[4]
        )

    def convert_template(self):
        with open(self.tempfile_path, "w+b") as output_file:
            pisa.CreatePDF(self.build_template(), dest=output_file, pagesize=(210, 203))

    def print_result(self):
        self.convert_template()
        os.startfile(self.tempfile_path)
