import tkinter as tk
from tkinter import ttk
import customtkinter
import os
import tempfile

from calculator.renderer.oil_renderer import OilRenderer
from calculator.renderer.briquettes_renderer import BriquettesRenderer


class CalcApp(customtkinter.CTk):

    WIDTH = 940
    HEIGHT = 520

    button_color = "#e67b36"
    font_color = "white"
    hover_color = "#D67334"
    sidebar_header_color = "#e67b36"
    header_color = "grey"
    formula_color = "white"
    copyright_color= "#6b6b6b"

    def __init__(self):
        super().__init__()

        self.delivery_quantity_briquettes = tk.StringVar(self)
        self.delivery_quantity_oil = tk.StringVar(self)
        self.result_briquettes = 0
        self.is_valid_calc = False

        # Window configuration
        self.title("Emissionsrechner")
        self.geometry(f"{CalcApp.WIDTH}x{CalcApp.HEIGHT}")
        self.minsize(self.WIDTH, self.HEIGHT)
        self.maxsize(self.WIDTH, self.HEIGHT)

        self.grid_columnconfigure(1, weight=1)
        self.grid_columnconfigure((2, 3), weight=0)
        self.grid_rowconfigure((0, 1, 2), weight=1)

        ## Sidebar 
        self.sidebar = customtkinter.CTkFrame(self, width=160, corner_radius=0, fg_color="#2b2b2b")
        self.sidebar.grid(row=0, column=0, rowspan=5, sticky="nsew")
        self.sidebar.grid_rowconfigure(5, weight=1)
        self.logo_label = customtkinter.CTkLabel(
            self.sidebar, 
            text="Emissionsrechner", 
            text_color=self.sidebar_header_color,
            font=customtkinter.CTkFont(size=30, weight="bold")
        )
        self.logo_label.grid(row=0, column=0, padx=20, pady=(20, 10))
        self.sidebar_briquettes = customtkinter.CTkButton(
            self.sidebar, 
            command=self.calc_briquettes, 
            text="Briketts",
            height=30,
            width=150,
            text_color=self.font_color,
            fg_color=self.button_color, 
            hover_color=self.hover_color,
            font=customtkinter.CTkFont(size=16)
        )
        self.sidebar_briquettes.grid(row=1, column=0, padx=20, pady=10)
        self.sidebar_heating_oil = customtkinter.CTkButton(
            self.sidebar, 
            command=self.calc_heating_oil, 
            text="Heizöl",
            height=30,
            width=150,
            text_color=self.font_color,
            fg_color=self.button_color, 
            hover_color=self.hover_color,
            font=customtkinter.CTkFont(size=16)
        )
        self.sidebar_heating_oil.grid(row=2, column= 0, padx=20, pady=10)
        self.sidebar_copyright = customtkinter.CTkLabel(
            self.sidebar,
            text="© Lycr.eu",
            font=customtkinter.CTkFont(size=16),
            text_color=self.copyright_color
        )
        self.sidebar_copyright.grid(row=7, column=0, padx=20, pady=10)

    ## functions
    def calc_briquettes(self):
        briquettes_renderer = BriquettesRenderer(self, self.button_color, self.font_color, self.hover_color, self.sidebar_header_color, self.header_color, self.formula_color)
        briquettes_renderer.build_base()
        
    def calc_heating_oil(self):
        oil_renderer = OilRenderer(self, self.button_color, self.font_color, self.hover_color, self.sidebar_header_color, self.header_color, self.formula_color)
        oil_renderer.build_base()

if __name__ == "__main__":
    app = CalcApp()
    app.mainloop()
    if os.path.isfile(os.path.join(tempfile.gettempdir(), 'label.pdf')):
        os.remove(os.path.join(tempfile.gettempdir(), 'label.pdf'))