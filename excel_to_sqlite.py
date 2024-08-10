import pandas as pd
from sqlalchemy import create_engine, Column, Integer, String
from sqlalchemy.orm import sessionmaker, declarative_base

Base = declarative_base()

# Map excel headers to table columns
column_mapping = {
    'Employee ID': 'eid',
    'First Name': 'fname',
    'Last Name': 'lname',
    'Address': 'addr'
}

class User(Base):
    __tablename__ = 'users'
    id = Column(Integer, primary_key=True)
    eid = Column(String, name=column_mapping['Employee ID'])
    fname = Column(String, name=column_mapping['First Name'])
    lname = Column(String, name=column_mapping['Last Name'])
    addr = Column(String, name=column_mapping['Address'])

def read_excel_to_sqlite(sheet_name, excel_file, db_file):
    """
    Reads data from an Excel file and writes it to a SQLite database.

    Parameters:
        sheet_name (str): The name of the Excel sheet to read from.
        excel_file (str): The path to the Excel file.
        db_file (str): The path to the SQLite database file.

    Returns:
        None
    """
    engine = create_engine(f'sqlite:///{db_file}')
    Base.metadata.create_all(engine)
    Session = sessionmaker(bind=engine)
    session = Session()

    df = pd.read_excel(excel_file, sheet_name=sheet_name)
    df = df[list(column_mapping.keys())]
    print("Read rows from Excel file:", df.shape[0])

    for index, row in df.iterrows():
        user = session.get(User, index)
        # Update user if it exists, otherwise create a new one
        if user and not all(getattr(user, column_mapping[col]) == val
                for col, val in row.items()):
            for col, val in row.items():
                setattr(user, column_mapping[col], val)
            session.merge(user)
        elif not user:
            new_user = User(
                id=index,
                **{column_mapping[col]: val for col, val in row.items()}
            )
            session.add(new_user)
    session.commit()
    session.close()

def main():
    """
    A function that reads data from an Excel file and writes it to a SQLite database.

    Parameters:
        None

    Returns:
        None
    """
    excel_file = 'users.xlsx'
    sheet_name = 'Sheet1_2'
    db_file = 'users.db'
    read_excel_to_sqlite(sheet_name, excel_file, db_file)

if __name__ == "__main__":
    main()