# Install PostgreSQL
### For Linux
sudo -u postgres initdb --locale en_US.UTF-8 -D /var/lib/postgres/data
sudo systemctl start postgresql

### Optional enable postgresql
sudo systemctl enable postgresql
